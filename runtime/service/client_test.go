package service

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/manifest"
)

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

func TestClient(t *testing.T) {
	addr := "/tmp/test-fold-ingress-server.sock"
	server := newTestIngressServer(addr)
	go func() {
		server.start(t)
	}()
	defer server.stop()
	client := newIngressClient(addr)
	err := client.start()
	if err != nil {
		t.Fatalf("Failed start client")
	}
	m, err := client.getManifest(context.Background())
	if err != nil {
		t.Fatalf("Failed to request manifest")
	}
	expectation := &manifest.Version{Major: 1, Minor: 0, Patch: 0}
	if !compareVersion(m.Version, expectation) {
		t.Fatalf("Exepcted manifest to have version %+v, but found %+v", expectation, m.Version)
	}

	req := &Request{HttpMethod: manifest.HttpMethod_GET, Path: "/test", Body: `{"msg": "test_body"}`, Headers: nil, Params: nil}
	res, err := client.doRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("Failed to make request")
	}
	if res.Body != req.Body {
		t.Fatalf("Exepcted respond body to equal request body. Expected %v but found %v", req.Body, res.Body)
	}
}
