package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/transport/pb"
	"github.com/foldsh/fold/runtime/types"
)

// Ingress wraps the gRPC client to communicate with the service.
// This is how the runtime manages inbound traffic for the service.
// It is down to the SDK to implement the server half of this spec.
type Ingress struct {
	foldSockAddr string
	conn         *grpc.ClientConn
	client       pb.FoldIngressClient
	logger       logging.Logger
}

// Creates a new `IngressClient`. The `foldSockAddr` should be a complete
// file path and it should match the one used to start the server.
func NewIngress(logger logging.Logger, foldSockAddr string) *Ingress {
	return &Ingress{foldSockAddr: foldSockAddr, logger: logger}
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func (ic *Ingress) Start() error {
	// We aren't bothering with a secure connection as it's all local
	// over a unix domain socket. We block to guarantee that by the time
	// the client is returned, the connection is alive and established.
	// This takes around 2 to 4 ms usually, and the backoff config ensures
	// that we return almost as soon as it's up. The default backoff
	// config waits for a second, which is pointless for us.
	ic.logger.Debugf("dialing server on %s", ic.foldSockAddr)
	conn, err := grpc.Dial(
		ic.foldSockAddr,
		grpc.WithInsecure(),
		grpc.WithAuthority("localhost"),
		grpc.WithDialer(dialer),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			backoff.Config{
				500 * time.Microsecond,
				1.1,
				0.2,
				4 * time.Second,
			},
			500 * time.Microsecond,
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to dial grpc server: %v", err)
	}
	ic.conn = conn
	ic.client = pb.NewFoldIngressClient(conn)
	return nil
}

// Retrieve the service manifest.
func (ic *Ingress) GetManifest(ctx context.Context) (*manifest.Manifest, error) {
	return ic.client.GetManifest(ctx, &pb.ManifestReq{})
}

// Submit a request to the service for processing.
func (ic *Ingress) DoRequest(
	ctx context.Context,
	in *types.Request,
) (*types.Response, error) {
	encoded, err := in.ToProto()
	if err != nil {
		return nil, err
	}
	res, err := ic.client.DoRequest(ctx, encoded)
	return types.ResFromProto(res), err
}
