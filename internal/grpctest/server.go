package grpctest

import (
	"context"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/transport/pb"
)

type Server struct {
	pb.UnimplementedFoldIngressServer
	socket   string
	server   *grpc.Server
	manifest *manifest.Manifest
	t        *testing.T
	logger   logging.Logger

	ManifestCalls  int
	DoRequestCalls int
	LastRequest    *manifest.FoldHTTPRequest
}

func NewServer(t *testing.T, logger logging.Logger, foldSockAddr string) *Server {
	manifest := &manifest.Manifest{
		Version: &manifest.Version{Major: 1, Minor: 0, Patch: 0},
	}
	return &Server{socket: foldSockAddr, manifest: manifest, t: t, logger: logger}
}

func (s *Server) Start() {
	lis, err := net.Listen("unix", s.socket)
	if err != nil {
		s.t.Fatalf("%+v", err)
	}
	s.server = grpc.NewServer()
	pb.RegisterFoldIngressServer(s.server, s)
	if err := s.server.Serve(lis); err != nil {
		s.t.Fatalf("%+v", err)
	}
}

func (s *Server) Stop() {
	os.Remove(s.socket)
	s.server.Stop()
}

func (s *Server) GetManifest(
	ctx context.Context,
	in *pb.ManifestReq,
) (*manifest.Manifest, error) {
	s.logger.Debugf("Handling GetManifest")
	s.ManifestCalls++
	return s.manifest, nil
}

func (s *Server) DoRequest(
	ctx context.Context,
	in *manifest.FoldHTTPRequest,
) (*manifest.FoldHTTPResponse, error) {
	s.logger.Debugf("Handling DoRequest")
	s.DoRequestCalls++
	s.LastRequest = in
	return &manifest.FoldHTTPResponse{Status: 200, Body: in.Body, Headers: nil}, nil
}
