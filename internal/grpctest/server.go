package grpctest

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/foldsh/fold/manifest"
)

type Ingress struct {
	pb.UnimplementedFoldIngress
	socket   string
	server   *grpc.Server
	manifest *manifest.Manifest
}

func NewIngress(foldSockAddr string) *Ingress {
	restartCount += 1
	manifest := &manifest.Manifest{
		Version: &manifest.Version{Major: int32(restartCount), Minor: 0, Patch: 0},
	}
	return &Ingress{socket: foldSockAddr, manifest: manifest}
}

func (is *Ingress) Start() {
	lis, err := net.Listen("unix", is.socket)
	if err != nil {
		panic(err)
	}
	is.server = grpc.NewServer()
	pb.RegisterFoldIngress(is.server, is)
	if err := is.server.Serve(lis); err != nil {
		panic(err)
	}
}

func (is *Ingress) Stop() {
	is.server.Stop()
}

func (is *Ingress) GetManifest(
	ctx context.Context,
	in *pb.ManifestReq,
) (*manifest.Manifest, error) {
	return is.manifest, nil
}

func (is *Ingress) DoRequest(
	ctx context.Context,
	in *manifest.FoldHTTPRequest,
) (*manifest.FoldHTTPResponse, error) {
	return &manifest.FoldHTTPResponse{Status: 200, Body: in.Body, Headers: nil}, nil
}

func compareVersion(a *manifest.Version, b *manifest.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch
}
