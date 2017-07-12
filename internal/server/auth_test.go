package server

import (
	"context"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"gitlab.com/gitlab-org/gitaly/auth"
	"gitlab.com/gitlab-org/gitaly/internal/config"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netctx "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const scratchDir = "testdata"

var serverSocketPath = path.Join(scratchDir, "gitaly.sock")

func TestSanity(t *testing.T) {
	srv := runServer(t)
	defer srv.Stop()
	connOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := dial(connOpts)
	require.NoError(t, err)
	defer conn.Close()

	require.NoError(t, healthCheck(conn))
}

func TestAuthFailures(t *testing.T) {
	defer func(oldAuth config.Auth) {
		config.Config.Auth = oldAuth
	}(config.Config.Auth)
	config.Config.Auth.Token = config.Token("quxbaz")

	srv := runServer(t)
	defer srv.Stop()

	testCases := []struct {
		desc string
		opts []grpc.DialOption
		code codes.Code
	}{
		{desc: "no auth", opts: nil, code: codes.Unauthenticated},
		{
			desc: "invalid auth",
			opts: []grpc.DialOption{grpc.WithPerRPCCredentials(brokenAuth{})},
			code: codes.Unauthenticated,
		},
		{
			desc: "wrong secret",
			opts: []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentials("foobar"))},
			code: codes.PermissionDenied,
		},
	}
	for _, tc := range testCases {
		t.Log(tc.desc)
		connOpts := append(tc.opts, grpc.WithInsecure())
		func() {
			conn, err := dial(connOpts)
			require.NoError(t, err, tc.desc)
			defer conn.Close()
			testhelper.AssertGrpcError(t, healthCheck(conn), tc.code, "")
		}()
	}
}

func TestAuthSuccess(t *testing.T) {
	defer func(oldAuth config.Auth) {
		config.Config.Auth = oldAuth
	}(config.Config.Auth)

	srv := runServer(t)
	defer srv.Stop()

	testCases := []struct {
		desc     string
		opts     []grpc.DialOption
		required bool
		token    config.Token
	}{
		{desc: "no auth, not required"},
		{
			desc:  "incorrect auth, not required",
			opts:  []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentials("incorrect"))},
			token: config.Token("foobar"),
		},
		{
			desc:  "correct auth, not required",
			opts:  []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentials("foobar"))},
			token: config.Token("foobar"),
		},
		{
			desc:     "correct auth, required",
			opts:     []grpc.DialOption{grpc.WithPerRPCCredentials(gitalyauth.RPCCredentials("foobar"))},
			token:    config.Token("foobar"),
			required: true,
		},
	}
	for _, tc := range testCases {
		config.Config.Auth.Token = tc.token
		config.Config.Auth.Transitioning = !tc.required
		t.Logf("%+v", config.Config.Auth)
		connOpts := append(tc.opts, grpc.WithInsecure())
		func() {
			conn, err := dial(connOpts)
			require.NoError(t, err, tc.desc)
			defer conn.Close()
			assert.NoError(t, healthCheck(conn), tc.desc)
		}()
	}
}

type brokenAuth struct{}

func (brokenAuth) RequireTransportSecurity() bool { return false }
func (brokenAuth) GetRequestMetadata(netctx.Context, ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer blablabla"}, nil
}

func dial(opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithDialer(func(addr string, _ time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}))
	return grpc.Dial(serverSocketPath, opts...)
}

func healthCheck(conn *grpc.ClientConn) error {
	client := healthpb.NewHealthClient(conn)
	_, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{})
	return err
}

func runServer(t *testing.T) *grpc.Server {
	srv := New()
	if err := os.Remove(serverSocketPath); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if err := os.MkdirAll(scratchDir, 0755); err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("unix", serverSocketPath)
	require.NoError(t, err)
	go srv.Serve(listener)
	return srv
}