package commit

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strconv"

	"gitlab.com/gitlab-org/gitaly/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/internal/helper"

	pb "gitlab.com/gitlab-org/gitaly-proto/go"
	"gitlab.com/gitlab-org/gitaly/streamio"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type entryInfo struct {
	objectType string
	oid        string
	mode       int32
}

func (s *server) TreeEntry(in *pb.TreeEntryRequest, stream pb.CommitService_TreeEntryServer) error {
	if err := validateRequest(in); err != nil {
		return grpc.Errorf(codes.InvalidArgument, "TreeEntry: %v", err)
	}

	repoPath, err := helper.GetRepoPath(in.Repository)
	if err != nil {
		return err
	}

	stdinReader, stdinWriter := io.Pipe()
	cmdArgs := []string{
		"--git-dir", repoPath,
		"cat-file",
		"--batch",
	}
	cmd, err := helper.NewCommand(exec.Command(helper.GitPath(), cmdArgs...), stdinReader, nil, nil)
	if err != nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: cmd: %v", err)
	}
	defer cmd.Kill()
	defer stdinWriter.Close()
	defer stdinReader.Close()

	dirName := path.Dir(string(in.GetPath()))
	if dirName == "." {
		dirName = ""
	}
	baseName := path.Base(string(in.GetPath()))

	treeObject := fmt.Sprintf("%s^{tree}:%s\n", in.GetRevision(), dirName)
	if _, err := stdinWriter.Write([]byte(treeObject)); err != nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: stdin write: %v", err)
	}

	stdout := bufio.NewReader(cmd)

	treeInfo, err := catfile.ParseObjectInfo(stdout)
	if err != nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: %v", err)
	}
	if treeInfo.Oid == "" {
		return sendNotFoundResponse(stream)
	}

	treeEntryInfo, err := extractEntryInfoFromTreeData(stdout, treeInfo.Size, baseName)
	if err != nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: %v", err)
	}
	if treeEntryInfo.oid == "" {
		return sendNotFoundResponse(stream)
	}

	if treeEntryInfo.objectType == "commit" {
		response := &pb.TreeEntryResponse{
			Type: pb.TreeEntryResponse_COMMIT,
			Mode: treeEntryInfo.mode,
			Oid:  treeEntryInfo.oid,
		}
		if err := stream.Send(response); err != nil {
			return grpc.Errorf(codes.Unavailable, "TreeEntry: send: %v", err)
		}

		return nil
	}

	stdinWriter.Write([]byte(treeEntryInfo.oid))
	stdinWriter.Close()

	objectInfo, err := catfile.ParseObjectInfo(stdout)
	if err != nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: %v", err)
	}

	if treeEntryInfo.objectType != objectInfo.Type {
		return grpc.Errorf(
			codes.Internal,
			"TreeEntry: mismatched object type: tree-oid=%s object-oid=%s entry-type=%s object-type=%s",
			treeEntryInfo.oid, objectInfo.Oid, treeEntryInfo.objectType, objectInfo.Type,
		)
	}

	if objectInfo.Type == "tree" {
		response := &pb.TreeEntryResponse{
			Type: pb.TreeEntryResponse_TREE,
			Oid:  objectInfo.Oid,
			Size: objectInfo.Size,
			Mode: treeEntryInfo.mode,
		}
		return helper.DecorateError(codes.Unavailable, stream.Send(response))
	}

	dataLength := objectInfo.Size
	if in.Limit > 0 && dataLength > in.Limit {
		dataLength = in.Limit
	}

	response := &pb.TreeEntryResponse{
		Type: pb.TreeEntryResponse_BLOB,
		Oid:  objectInfo.Oid,
		Size: objectInfo.Size,
		Mode: treeEntryInfo.mode,
	}
	if dataLength == 0 {
		return helper.DecorateError(codes.Unavailable, stream.Send(response))
	}

	sw := streamio.NewWriter(func(p []byte) error {
		response.Data = p

		if err := stream.Send(response); err != nil {
			return grpc.Errorf(codes.Unavailable, "TreeEntry: send: %v", err)
		}

		// Use a new response so we don't send other fields (Size, ...) over and over
		response = &pb.TreeEntryResponse{}

		return nil
	})

	n, err := io.Copy(sw, io.LimitReader(stdout, dataLength))
	if n < dataLength && err == nil {
		return grpc.Errorf(codes.Internal, "TreeEntry: Incomplete copy")
	}

	return err
}

func validateRequest(in *pb.TreeEntryRequest) error {
	if len(in.GetRevision()) == 0 {
		return fmt.Errorf("empty Revision")
	}

	if len(in.GetPath()) == 0 {
		return fmt.Errorf("empty Path")
	}

	return nil
}

func extractEntryInfoFromTreeData(stdout *bufio.Reader, treeSize int64, baseName string) (*entryInfo, error) {
	var modeBytes, path []byte
	var objectType string
	var err error

	oidBytes := make([]byte, 20)
	entryFound := false
	bytesLeft := treeSize

	for bytesLeft > 0 {
		modeBytes, err = stdout.ReadBytes(' ')
		if err != nil || len(modeBytes) <= 1 {
			return nil, fmt.Errorf("read entry mode: %v", err)
		}
		bytesLeft -= int64(len(modeBytes))
		modeBytes = modeBytes[:len(modeBytes)-1]

		path, err = stdout.ReadBytes('\x00')
		if err != nil || len(path) <= 1 {
			return nil, fmt.Errorf("read entry path: %v", err)
		}
		bytesLeft -= int64(len(path))
		path = path[:len(path)-1]

		if n, _ := stdout.Read(oidBytes); n != 20 {
			return nil, fmt.Errorf("read entry oid: %v", err)
		}

		bytesLeft -= int64(len(oidBytes))

		if string(path) == baseName {
			entryFound = true
			break
		}
	}

	// Extra byte for a linefeed at the end
	if _, err := stdout.Discard(int(bytesLeft + 1)); err != nil {
		return nil, fmt.Errorf("stdout discard: %v", err)
	}

	if !entryFound {
		return &entryInfo{}, nil
	}

	mode, err := strconv.ParseInt(string(modeBytes), 8, 32)
	if err != nil {
		return nil, fmt.Errorf("parse mode: %v", err)
	}

	oid := fmt.Sprintf("%02x", oidBytes)

	// Based on https://github.com/git/git/blob/v2.13.1/builtin/ls-tree.c#L67-L87
	switch mode & 0xf000 {
	case 0160000:
		objectType = "commit"
	case 040000:
		objectType = "tree"
	default:
		objectType = "blob"
	}

	return &entryInfo{
		objectType: objectType,
		mode:       int32(mode),
		oid:        oid,
	}, nil
}

func sendNotFoundResponse(stream pb.Commit_TreeEntryServer) error {
	return helper.DecorateError(codes.Unavailable, stream.Send(&pb.TreeEntryResponse{}))
}