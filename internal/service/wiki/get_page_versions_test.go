package wiki

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"gitlab.com/gitlab-org/gitaly/internal/testhelper"

	pb "gitlab.com/gitlab-org/gitaly-proto/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWikiGetPageVersionsRequest(t *testing.T) {
	wikiRepo, wikiRepoPath, cleanupFunc := setupWikiRepo(t)
	defer cleanupFunc()

	ctx, cancel := testhelper.Context()
	defer cancel()

	server, serverSocketPath := runWikiServiceServer(t)
	defer server.Stop()

	client, conn := newWikiClient(t, serverSocketPath)
	defer conn.Close()
	pageTitle := "WikiGétPageVersions"

	content := bytes.Repeat([]byte("Mock wiki page content"), 10000)
	writeWikiPage(t, client, wikiRepo, createWikiPageOpts{title: pageTitle, content: content})
	v1cid := testhelper.MustRunCommand(t, nil, "git", "-C", wikiRepoPath, "log", "-1", "--format=%H")
	updateWikiPage(t, client, wikiRepo, pageTitle, []byte("New content"))
	v2cid := testhelper.MustRunCommand(t, nil, "git", "-C", wikiRepoPath, "log", "-1", "--format=%H")

	gitAuthor := &pb.CommitAuthor{
		Name:  []byte("Ahmad Sherif"),
		Email: []byte("ahmad@gitlab.com"),
	}

	testCases := []struct {
		desc     string
		request  *pb.WikiGetPageVersionsRequest
		versions []*pb.WikiPageVersion
	}{
		{
			desc: "No page found",
			request: &pb.WikiGetPageVersionsRequest{
				Repository: wikiRepo,
				PagePath:   []byte("not-found"),
			},
			versions: nil,
		},
		{
			desc: "2 version found",
			request: &pb.WikiGetPageVersionsRequest{
				Repository: wikiRepo,
				PagePath:   []byte(pageTitle),
			},
			versions: []*pb.WikiPageVersion{
				{
					Commit: &pb.GitCommit{
						Id:        strings.TrimRight(string(v2cid), "\n"),
						Body:      []byte("Update WikiGétPageVersions"),
						Subject:   []byte("Update WikiGétPageVersions"),
						Author:    gitAuthor,
						Committer: gitAuthor,
						ParentIds: []string{strings.TrimRight(string(v1cid), "\n")},
					},
					Format: "markdown",
				},
				{
					Commit: &pb.GitCommit{
						Id:        strings.TrimRight(string(v1cid), "\n"),
						Body:      []byte("Add WikiGétPageVersions"),
						Subject:   []byte("Add WikiGétPageVersions"),
						Author:    gitAuthor,
						Committer: gitAuthor,
					},
					Format: "markdown",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			stream, err := client.WikiGetPageVersions(ctx, tc.request)
			require.NoError(t, err)
			require.NoError(t, stream.CloseSend())

			response, err := stream.Recv()
			require.NoError(t, err)

			require.Len(t, response.GetVersions(), len(tc.versions))
			for i, version := range response.GetVersions() {
				v2 := tc.versions[i]
				assertWikiPageVersionEqual(t, version, v2, "%d blew up", i)
			}
		})
	}
}

func TestWikiGetPageVersionsPaginationParams(t *testing.T) {
	wikiRepo, _, cleanupFunc := setupWikiRepo(t)
	defer cleanupFunc()

	ctx, cancel := testhelper.Context()
	defer cancel()

	server, serverSocketPath := runWikiServiceServer(t)
	defer server.Stop()

	client, conn := newWikiClient(t, serverSocketPath)
	defer conn.Close()

	pageTitle := "WikiGetPageVersions"
	content := []byte("page content")
	writeWikiPage(t, client, wikiRepo, createWikiPageOpts{title: pageTitle, content: content})

	for i := 0; i < 25; i++ {
		updateWikiPage(t, client, wikiRepo, pageTitle, []byte(string(i)))
	}

	testCases := []struct {
		desc        string
		perPage     int32
		page        int32
		nrOfResults int
	}{
		{
			desc:        "default to page 1 with 20 items",
			nrOfResults: 20,
		},
		{
			desc:        "oversized perPage param",
			perPage:     100,
			nrOfResults: 26,
		},
		{
			desc:        "allows later pages",
			page:        2,
			nrOfResults: 6,
		},
		{
			desc:        "returns nothing of no versions can be found",
			page:        100,
			nrOfResults: 0,
		},
		{
			// https://github.com/gollum/gollum-lib/blob/be6409315f6af5a6d90eb012a1154b485579db67/lib/gollum-lib/pagination.rb#L34
			desc:        "per page is minimal 20",
			perPage:     1,
			nrOfResults: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			stream, err := client.WikiGetPageVersions(ctx, &pb.WikiGetPageVersionsRequest{
				Repository: wikiRepo,
				PagePath:   []byte(pageTitle),
				PerPage:    tc.perPage,
				Page:       tc.page})
			require.NoError(t, err)

			nrFoundVersions := 0
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					t.Fatal(err)
				}

				nrFoundVersions += len(resp.GetVersions())
			}

			require.Equal(t, tc.nrOfResults, nrFoundVersions)
		})
	}
}

func assertWikiPageVersionEqual(t *testing.T, expected, actual *pb.WikiPageVersion, msg ...interface{}) bool {
	assert.Equal(t, expected.GetFormat(), actual.GetFormat())
	assert.NoError(t, testhelper.GitCommitEqual(expected.GetCommit(), actual.GetCommit()))

	switch len(msg) {
	case 0:
		t.Logf("WikiPageVersion differs: %v != %v", expected, actual)
	case 1:
		t.Log(msg[0])
	default:
		t.Logf(msg[0].(string), msg[1:])
	}

	return t.Failed()
}
