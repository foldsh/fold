package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/transport/pb"
)

// Ingress wraps the gRPC client to communicate with the service.
// This is how the runtime manages inbound traffic for the service.
// It is down to the SDK to implement the server half of this spec.
type Ingress struct {
	conn   *grpc.ClientConn
	client pb.FoldIngressClient
	logger logging.Logger
}

// Creates a new `IngressClient`. The `foldSockAddr` should be a complete
// file path and it should match the one used to start the server.
func NewIngress(logger logging.Logger) *Ingress {
	return &Ingress{logger: logger}
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func (i *Ingress) Start(socketAddress string) error {
	// We aren't bothering with a secure connection as it's all local
	// over a unix domain socket. We block to guarantee that by the time
	// the client is returned, the connection is alive and established.
	// This takes around 2 to 4 ms usually, and the backoff config ensures
	// that we return almost as soon as it's up. The default backoff
	// config waits for a second, which is pointless for us.
	i.logger.Debugf("Dialing server on %s", socketAddress)
	conn, err := grpc.Dial(
		socketAddress,
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
	i.conn = conn
	i.client = pb.NewFoldIngressClient(conn)
	i.logger.Debugf("Connected")
	return nil
}

func (i *Ingress) Stop() error {
	if err := i.conn.Close(); err != nil {
		if grpc.Code(err) == codes.Canceled {
			return nil
		}
		return errors.New("failed to close the connection")
	}
	return nil
}

func (i *Ingress) Restart(socketAddress string) error {
	if err := i.Stop(); err != nil {
		return err
	}
	return i.Start(socketAddress)
}

// Retrieve the service manifest.
func (i *Ingress) GetManifest(ctx context.Context) (*manifest.Manifest, error) {
	return i.client.GetManifest(ctx, &pb.ManifestReq{})
}

// Submit a request to the service for processing.
func (i *Ingress) DoRequest(
	ctx context.Context,
	in *Request,
) (*Response, error) {
	encoded, err := in.ToProto()
	if err != nil {
		return nil, err
	}
	res, err := i.client.DoRequest(ctx, encoded)
	if res != nil && err == nil {
		return ResFromProto(res), err
	}
	return nil, err
}
