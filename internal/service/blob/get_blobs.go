package blob

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"gitlab.com/gitlab-org/gitaly/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/internal/service/commit"
	"gitlab.com/gitlab-org/gitaly/streamio"

	pb "gitlab.com/gitlab-org/gitaly-proto/go"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getBlobsHandler(req *pb.GetBlobsRequest, stream pb.BlobService_GetBlobsServer) catfile.Handler {
	return func(stdin io.Writer, stdout *bufio.Reader) error {
		for _, revisionPath := range req.RevisionPaths {
			revision := revisionPath.Revision
			path := revisionPath.Path

			treeEntry, err := commit.TreeEntryForRevisionAndPath(revision, string(path), stdin, stdout)
			if err != nil {
				return err
			}

			response := &pb.GetBlobsResponse{Revision: revision, Path: path}

			if treeEntry == nil || len(treeEntry.Oid) == 0 {
				if err := stream.Send(response); err != nil {
					return status.Errorf(codes.Unavailable, "GetBlobs: send: %v", err)
				}

				continue
			}

			response.Mode = treeEntry.Mode
			response.Oid = treeEntry.Oid

			if treeEntry.Type == pb.TreeEntry_COMMIT {
				response.IsSubmodule = true

				if err := stream.Send(response); err != nil {
					return status.Errorf(codes.Unavailable, "GetBlobs: send: %v", err)
				}

				continue
			}

			fmt.Fprintln(stdin, treeEntry.Oid)

			objectInfo, err := catfile.ParseObjectInfo(stdout)
			if err != nil {
				return status.Errorf(codes.Internal, "GetBlobs: %v", err)
			}
			if objectInfo.Type != "blob" {
				return status.Errorf(codes.InvalidArgument, "GetBlobs: object at %s:%s is %s, not blob", revision, path, objectInfo.Type)
			}

			response.Size = objectInfo.Size

			var readLimit int64
			if req.Limit < 0 || req.Limit > objectInfo.Size {
				readLimit = objectInfo.Size
			} else {
				readLimit = req.Limit
			}

			if readLimit == 0 {
				if err := stream.Send(response); err != nil {
					return status.Errorf(codes.Unavailable, "GetBlobs: send: %v", err)
				}
			} else {
				sw := streamio.NewWriter(func(p []byte) error {
					msg := &pb.GetBlobsResponse{}
					if response != nil {
						msg = response
						response = nil
					}

					msg.Data = p

					return stream.Send(msg)
				})

				sent, err := io.Copy(sw, io.LimitReader(stdout, readLimit))
				if err != nil {
					return status.Errorf(codes.Unavailable, "GetBlobs: send: %v", err)
				}
				if sent != readLimit {
					return status.Errorf(codes.Unavailable, "GetBlobs: short send: %d/%d bytes", sent, objectInfo.Size)
				}
			}

			discardSize := int64(1) // A new line after the blob content
			if readLimit < objectInfo.Size {
				discardSize += objectInfo.Size - readLimit
			}

			discarded, err := io.Copy(ioutil.Discard, io.LimitReader(stdout, discardSize))
			if err != nil {
				return status.Errorf(codes.Unavailable, "GetBlobs: discarding data: %v", err)
			}
			if discarded != discardSize {
				return status.Errorf(codes.Unavailable, "GetBlobs: short discard: %d/%d bytes", discarded, discardSize)
			}
		}

		return nil
	}
}

func (*server) GetBlobs(req *pb.GetBlobsRequest, stream pb.BlobService_GetBlobsServer) error {
	if req.Repository == nil {
		return status.Errorf(codes.InvalidArgument, "GetBlobs: empty Repository")
	}

	if len(req.RevisionPaths) == 0 {
		return status.Errorf(codes.InvalidArgument, "GetBlobs: empty RevisionPaths")
	}

	handler := getBlobsHandler(req, stream)
	return catfile.CatFile(stream.Context(), req.Repository, handler)
}
