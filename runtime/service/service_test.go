package service

import (
	"context"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

// This tests the whole lifecycle of a service, albeit in a fairly abstract
// way. For the sake of making the test easier I've made an implementation
// of the foldSubprocess that just uses a goroutine to simulate the interface.
// This means I can just make a little mock gRPC server in the test file and
// run it in a goroutine.
func TestServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	service := newTestService(t)
	err := service.Start()
	if err != nil {
		t.Fatalf("failed to start service")
	}

	m, err := service.GetManifest()
	if err != nil {
		t.Fatalf("Failed to request manifest")
	}
	expectation := &manifest.Version{Major: 1, Minor: 0, Patch: 0}
	if !compareVersion(m.Version, expectation) {
		t.Fatalf("Exepcted manifest to have version %+v, but found %+v", expectation, m.Version)
	}

	req := &Request{HttpMethod: manifest.HttpMethod_GET, Path: "/test", Body: `{"msg": "test_body"}`, Headers: nil, Params: nil}
	res, err := service.DoRequest(req)
	if err != nil {
		t.Fatalf("Failed to make request")
	}
	if res.Body != req.Body {
		t.Fatalf("Exepcted respond body to equal request body. Expected %v but found %v", req.Body, res.Body)
	}
	service.Stop()
}

func newTestService(t *testing.T) Service {
	addr := newAddr()
	client := newIngressClient(addr)
	process := &goSubprocess{newTestIngressServer(addr), t}
	return &service{Command{}, addr, client, process, logging.NewTestLogger()}
}

// goroutine based implementation of the foldSubprocess
type goSubprocess struct {
	server *testIngressServer
	t      *testing.T
}

func (gsp *goSubprocess) run() error {
	go func() {
		// A sleep better simulates a new process starting and makes
		// sure that our logic for waiting for the server to come up
		// is working properly.
		time.Sleep(42 * time.Millisecond)
		gsp.server.start(gsp.t)
	}()
	return nil
}

func (gsp *goSubprocess) wait() error {
	return nil
}

func (gsp *goSubprocess) kill() error {
	gsp.server.stop()
	return nil
}

func (gsp *goSubprocess) signal(sig os.Signal) error {
	return nil
}

func (gsp *goSubprocess) setStdout(w io.Writer) {
	return
}

func (gsp *goSubprocess) setStderr(w io.Writer) {
	return
}

// ingressServer for testing
type testIngressServer struct {
	UnimplementedFoldIngressServer
	socket string
	server *grpc.Server
}

func newTestIngressServer(foldSockAddr string) *testIngressServer {
	return &testIngressServer{socket: foldSockAddr}
}

func (is *testIngressServer) start(t *testing.T) {
	lis, err := net.Listen("unix", is.socket)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	is.server = grpc.NewServer()
	RegisterFoldIngressServer(is.server, is)
	if err := is.server.Serve(lis); err != nil {
		t.Fatalf("failed to serve: %v", err)
	}
}

func (is *testIngressServer) stop() {
	is.server.Stop()
}

func (is *testIngressServer) GetManifest(ctx context.Context, in *ManifestReq) (*manifest.Manifest, error) {
	return &manifest.Manifest{Version: &manifest.Version{Major: 1, Minor: 0, Patch: 0}}, nil
}

func (is *testIngressServer) DoRequest(ctx context.Context, in *Request) (*Response, error) {
	return &Response{Status: 200, Body: in.Body, Headers: nil}, nil
}

func compareVersion(a *manifest.Version, b *manifest.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch
}
