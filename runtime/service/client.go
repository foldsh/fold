package service

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/foldsh/fold/manifest"
)

// ingressClient wraps the gRPC client to communicate with the service.
// This is how the runtime manages inbound traffic for the service.
// It is down to the SDK to implement the server half of this spec.
type ingressClient struct {
	foldSockAddr string
	conn         *grpc.ClientConn
	client       FoldIngressClient
}

// Creates a new `IngressClient`. The `foldSockAddr` should be a complete
// file path and it should match the one used to start the server.
func newIngressClient(foldSockAddr string) *ingressClient {
	return &ingressClient{foldSockAddr, nil, nil}
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func (ic *ingressClient) start() error {
	// We aren't bothering with a secure connection as it's all local
	// over a unix domain socket. We block to guarantee that by the time
	// the client is returned, the connection is alive and established.
	// This takes around 2 to 4 ms usually, and the backoff config ensures
	// that we return almost as soon as it's up. The default backoff
	// config waits for a second, which is pointless for us.
	conn, err := grpc.Dial(
		ic.foldSockAddr, grpc.WithInsecure(),
		grpc.WithDialer(dialer),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			backoff.Config{
				500 * time.Microsecond,
				1.6,
				0.2,
				16 * time.Second,
			},
			500 * time.Microsecond,
		}),
	)
	if err != nil {
		return fmt.Errorf("Failed to dial grpc server: %v", err)
	}
	ic.conn = conn
	ic.client = NewFoldIngressClient(conn)
	return nil
}

// Retrieve the service manifest.
func (ic *ingressClient) getManifest(ctx context.Context) (*manifest.Manifest, error) {
	return ic.client.GetManifest(ctx, &ManifestReq{})
}

// Submit a request to the service for processing.
func (ic *ingressClient) doRequest(ctx context.Context, in *Request) (*Response, error) {
	return ic.client.DoRequest(ctx, in)
}
