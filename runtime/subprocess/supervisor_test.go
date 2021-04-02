package supervisor

import (
	"context"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/subprocess/pb"
	"github.com/foldsh/fold/runtime/transport"
)

var restartCount = 0

// This tests the whole lifecycle of a service, albeit in a fairly abstract
// way. For the sake of making the test easier I've made an implementation
// of the foldSubprocess that just uses a goroutine to simulate the interface.
// This means I can just make a little mock gRPC server in the test file and
// run it in a goroutine.
func TestSupervisorIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	service := newTestSupervisor(t)
	service.Start()

	m, err := service.GetManifest()
	if err != nil {
		t.Fatalf("Failed to request manifest")
	}
	expectation := &manifest.Version{Major: 1, Minor: 0, Patch: 0}
	if !compareVersion(m.Version, expectation) {
		t.Fatalf("Exepcted manifest to have version %+v, but found %+v", expectation, m.Version)
	}

	req := &transport.Request{
		HTTPMethod: "GET",
		Body:       []byte(`{"msg": "test_body"}`),
		Route:      "/test",
	}
	res, err := service.DoRequest(req)
	if err != nil {
		t.Fatalf("Failed to make request")
	}
	if string(res.Body) != string(req.Body) {
		t.Fatalf(
			"Exepcted respond body to equal request body. Expected %v but found %v",
			req.Body,
			res.Body,
		)
	}
	// Now, if we restart the service, we should expect to see version two.
	service.Restart()
	m, err = service.GetManifest()
	if err != nil {
		t.Fatalf("Failed to request manifest")
	}
	expectation = &manifest.Version{Major: 2, Minor: 0, Patch: 0}
	if !compareVersion(m.Version, expectation) {
		t.Fatalf("Exepcted manifest to have version %+v, but found %+v", expectation, m.Version)
	}

	// When we shut down there should be no error.
	if err := service.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("+%v", err)
	}
}

func TestSupervisorShouldFunctionWithDeadProcessIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	service := newTestSupervisor(t)
	service.Start()
	service.process.kill()
	req := &transport.Request{
		HTTPMethod: "GET",
		Body:       []byte(`{"msg": "test_body"}`),
		Route:      "/test",
	}
	if res, err := service.DoRequest(req); err != nil {
		t.Fatalf("Failed to make request")
	} else {
		service.logger.Debugf("%+v", res)
		if res.Status != 502 {
			t.Errorf("Should have returned 502 with dead process")
		}
		testutils.Diff(
			t, string(CannotServiceRequest),
			string(res.Body),
			"Expected body to match CannotServiceRequest message",
		)
	}
}

func newTestSupervisor(t *testing.T) *Supervisor {
	logger := logging.NewTestLogger()
	return &Supervisor{
		clientFactory:     newIngressClient,
		subprocessFactory: newGoSubprocess,
		logger:            logger,
	}
}

// goroutine based implementation of the foldSubprocess
type goSubprocess struct {
	server *testIngressServer
	health bool
}

func newGoSubprocess(_ logging.Logger, addr string) foldSubprocess {
	return &goSubprocess{newTestIngressServer(addr), true}
}

func (gsp *goSubprocess) run(_ string, _ ...string) error {
	go func() {
		// A sleep better simulates a new process starting and makes
		// sure that our logic for waiting for the server to come up
		// is working properly.
		time.Sleep(42 * time.Millisecond)
		gsp.server.start()
	}()
	return nil
}

func (gsp *goSubprocess) wait() error {
	gsp.kill()
	return nil
}

func (gsp *goSubprocess) kill() error {
	gsp.server.stop()
	gsp.health = false
	return nil
}

func (gsp *goSubprocess) signal(sig os.Signal) error {
	gsp.kill()
	return nil
}

func (gsp *goSubprocess) healthz() bool {
	return gsp.health
}

// ingressServer for testing
type testIngressServer struct {
	pb.UnimplementedFoldIngressServer
	socket   string
	server   *grpc.Server
	manifest *manifest.Manifest
}

func newTestIngressServer(foldSockAddr string) *testIngressServer {
	restartCount += 1
	manifest := &manifest.Manifest{
		Version: &manifest.Version{Major: int32(restartCount), Minor: 0, Patch: 0},
	}
	return &testIngressServer{socket: foldSockAddr, manifest: manifest}
}

func (is *testIngressServer) start() {
	lis, err := net.Listen("unix", is.socket)
	if err != nil {
		panic(err)
	}
	is.server = grpc.NewServer()
	pb.RegisterFoldIngressServer(is.server, is)
	if err := is.server.Serve(lis); err != nil {
		panic(err)
	}
}

func (is *testIngressServer) stop() {
	is.server.Stop()
}

func (is *testIngressServer) GetManifest(
	ctx context.Context,
	in *pb.ManifestReq,
) (*manifest.Manifest, error) {
	return is.manifest, nil
}

func (is *testIngressServer) DoRequest(
	ctx context.Context,
	in *manifest.FoldHTTPRequest,
) (*manifest.FoldHTTPResponse, error) {
	return &manifest.FoldHTTPResponse{Status: 200, Body: in.Body, Headers: nil}, nil
}

func compareVersion(a *manifest.Version, b *manifest.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch
}
