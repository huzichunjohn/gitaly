package repository

import (
	"path"
	"testing"

	"gitlab.com/gitlab-org/gitaly/internal/testhelper"

	pb "gitlab.com/gitlab-org/gitaly-proto/go"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestSuccessfulIsSquashInProgressRequest(t *testing.T) {
	server, serverSocketPath := runRepoServer(t)
	defer server.Stop()

	client, conn := newRepositoryClient(t, serverSocketPath)
	defer conn.Close()

	testRepo1, testRepo1Path, cleanupFn := testhelper.NewTestRepo(t)
	defer cleanupFn()

	testhelper.MustRunCommand(t, nil, "git", "-C", testRepo1Path, "worktree", "add", "--detach", path.Join(testRepo1Path, "gitlab-worktree", "squash-1"), "master")

	testRepo2, _, cleanupFn := testhelper.NewTestRepo(t)
	defer cleanupFn()

	testCases := []struct {
		desc       string
		request    *pb.IsSquashInProgressRequest
		inProgress bool
	}{
		{
			desc: "Squash in progress",
			request: &pb.IsSquashInProgressRequest{
				Repository: testRepo1,
				SquashId:   "1",
			},
			inProgress: true,
		},
		{
			desc: "no Squash in progress",
			request: &pb.IsSquashInProgressRequest{
				Repository: testRepo2,
				SquashId:   "2",
			},
			inProgress: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()

			response, err := client.IsSquashInProgress(ctx, testCase.request)
			require.NoError(t, err)

			require.Equal(t, testCase.inProgress, response.InProgress)
		})
	}
}

func TestFailedIsSquashInProgressRequestDueToValidations(t *testing.T) {
	server, serverSocketPath := runRepoServer(t)
	defer server.Stop()

	client, conn := newRepositoryClient(t, serverSocketPath)
	defer conn.Close()

	testCases := []struct {
		desc    string
		request *pb.IsSquashInProgressRequest
		code    codes.Code
	}{
		{
			desc:    "empty repository",
			request: &pb.IsSquashInProgressRequest{SquashId: "1"},
			code:    codes.InvalidArgument,
		},
		{
			desc:    "empty Squash id",
			request: &pb.IsSquashInProgressRequest{Repository: &pb.Repository{}},
			code:    codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()

			_, err := client.IsSquashInProgress(ctx, testCase.request)
			testhelper.AssertGrpcError(t, err, testCase.code, "")
		})
	}
}
