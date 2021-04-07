package transport_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/foldsh/fold/internal/grpctest"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/transport"
)

func TestIngressCanDoRequestAndGetManifest(t *testing.T) {
	addr := "/tmp/fold.client.test-do-request-and-get-manifest.sock"
	client, server, _ := makeIngress(t, addr, 100*time.Millisecond)
	defer server.Stop()
	if err := client.Start(addr); err != nil {
		t.Fatalf("%+v", err)
	}

	client.GetManifest(context.Background())
	if server.ManifestCalls != 1 {
		t.Fatalf("Expected to record 1 manifest call but found %d", server.ManifestCalls)
	}

	req := &transport.Request{HTTPMethod: "GET", Body: []byte(`fold`)}
	client.DoRequest(context.Background(), req)
	if server.DoRequestCalls != 1 {
		t.Fatalf("Expected to record 1 do request call but found %d", server.DoRequestCalls)
	}
	if string(server.LastRequest.Body) != string(req.Body) {
		t.Fatalf(
			"Expected the body of the recorded gRPC request to have the body passed to the client.",
		)
	}
}

func TestIngressStop(t *testing.T) {
	addr := "/tmp/fold.client.test-stop.sock"
	client, server, _ := makeIngress(t, addr, 0)
	defer server.Stop()
	if err := client.Start(addr); err != nil {
		t.Fatalf("%+v", err)
	}
	if err := client.Stop(); err != nil {
		t.Fatalf("%+v", err)
	}
	req := &transport.Request{HTTPMethod: "GET", Body: []byte(`fold`)}
	_, err := client.DoRequest(context.Background(), req)
	if err == nil {
		t.Fatalf("Exepcted an error after stopping the client but no error was found")
	}
	if grpc.Code(err) != codes.Canceled {
		t.Fatalf("Expected an error code of Cancelled but found %v", grpc.Code(err))
	}
}

func TestIngressStopIsIdempotent(t *testing.T) {
	addr := "/tmp/fold.client.test-stop.sock"
	client, server, _ := makeIngress(t, addr, 0)
	defer server.Stop()
	if err := client.Start(addr); err != nil {
		t.Fatalf("%+v", err)
	}
	if err := client.Stop(); err != nil {
		t.Fatalf("%+v", err)
	}
	if err := client.Stop(); err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestIngressRestart(t *testing.T) {
	addr := "/tmp/fold.client.test-restart.sock"
	client, server, logger := makeIngress(t, addr, 0)
	defer server.Stop()
	if err := client.Start(addr); err != nil {
		t.Fatalf("%+v", err)
	}

	// The implementation of this test is a little convoluted to be honest, I should probably just
	// have mocked the connection.  Anyway, basically we will kick off a goroutine that will
	// endlessly make requests with the client. While that's running, we will then restart the
	// client under it. If the client actually restarts then we would expect to see the requests
	// succeeding, then failing, then start succeeding again.
	var (
		successBeforeRestart int
		failureDuringRestart int
		successAfterRestart  int
	)

	// We've got a potentially infinite loop in there, so a timeout is important.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// This channel is used to make sure that we get at least one request off before the restart.
	ready := make(chan struct{}, 1)
	// This channel is used to stop the goroutine and indicate that it has completed its sequence.
	done := make(chan struct{}, 1)

	go func() {
		var (
			started     bool
			failureSeen bool
		)
		for {
			select {
			case <-ctx.Done():
				logger.Debugf("Timed out")
				done <- struct{}{}
				return
			default:
				logger.Debugf("Running request")
				req := &transport.Request{HTTPMethod: "GET", Body: []byte(`fold`)}
				_, err := client.DoRequest(context.Background(), req)
				if err == nil && !failureSeen {
					successBeforeRestart++
				}
				if err == nil && failureSeen {
					successAfterRestart++
					done <- struct{}{}
					return
				}
				if err != nil {
					failureDuringRestart++
					failureSeen = true
				}
				if !started {
					ready <- struct{}{}
					started = true
				}
			}
		}
	}()

	<-ready

	if err := client.Restart(addr); err != nil {
		t.Fatalf("%+v", err)
	}

	<-done

	if !(successBeforeRestart > 0) {
		t.Errorf("Expected at least one success before the restart")
	}
	if !(failureDuringRestart > 0) {
		t.Errorf("Expected at least one failure during the restart")
	}
	if !(successAfterRestart > 0) {
		t.Errorf("Expected at least one success after the restart")
	}
}

func makeIngress(
	t *testing.T,
	addr string,
	delay time.Duration,
) (*transport.Ingress, *grpctest.Server, logging.Logger) {
	logger := logging.NewTestLogger()
	client := transport.NewIngress(logger)
	server := grpctest.NewServer(t, logger, addr)
	go func() {
		// A sleep makes sure that the client waits appropriately for the server to come up.
		time.Sleep(delay)
		server.Start()
	}()
	return client, server, logger
}
