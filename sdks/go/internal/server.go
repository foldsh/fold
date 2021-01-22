package internal

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/supervisor/pb"
)

type GrpcServer struct {
	pb.UnimplementedFoldIngressServer
	server *grpc.Server
}

func (s *GrpcServer) Start() {
	foldSockAddr := os.Getenv("FOLD_SOCK_ADDR")
	lis, err := net.Listen("unix", foldSockAddr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s.server = grpc.NewServer()
	pb.RegisterFoldIngressServer(s.server, s)
	if err := is.server.Serve(lis); err != nil {
		t.Fatalf("failed to serve: %v", err)
	}
}

func (s *GrpcServer) GetManifest(
	ctx context.Context,
	in *pb.ManifestReq,
) (*manifest.Manifest, error) {
	return &manifest.Manifest{Version: &manifest.Version{Major: 1, Minor: 0, Patch: 0}}, nil
}

func (s *GrpcServer) DoRequest(
	ctx context.Context,
	in *pb.Request,
) (*pb.Response, error) {
	return &pb.Response{Status: 200, Body: in.Body, Headers: nil}, nil
}
