package transport_test

import (
	"context"
	"testing"
	"time"

	"github.com/foldsh/fold/internal/grpctest"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/transport"
	"github.com/foldsh/fold/runtime/types"
)

func TestIngress(t *testing.T) {
	addr := "/tmp/fold.client.test.sock"
	logger := logging.NewTestLogger()
	client := transport.NewIngress(logger, addr)
	server := grpctest.NewServer(t, logger, addr)
	go func() {
		// A sleep makes sure that the client waits appropriately for the server to come up.
		time.Sleep(100 * time.Millisecond)
		server.Start()
	}()
	defer server.Stop()
	if err := client.Start(); err != nil {
		t.Fatalf("%+v", err)
	}

	client.GetManifest(context.Background())
	if server.ManifestCalls != 1 {
		t.Fatalf("Expected to record 1 manifest call but found %d", server.ManifestCalls)
	}

	req := &types.Request{HTTPMethod: "GET", Body: []byte(`fold`)}
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
