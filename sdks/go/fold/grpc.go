package fold

import (
	"context"
	"encoding/json"
	"net"
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/supervisor/pb"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedFoldIngressServer
	server  *grpc.Server
	service *service
	logger  logging.Logger
}

func (gs *grpcServer) start() {
	foldSockAddr := os.Getenv("FOLD_SOCK_ADDR")
	lis, err := net.Listen("unix", foldSockAddr)
	if err != nil {
		gs.logger.Fatalf("gRPC server failed to listen: %v", err)
	}
	gs.server = grpc.NewServer()
	pb.RegisterFoldIngressServer(gs.server, gs)
	if err := gs.server.Serve(lis); err != nil {
		gs.logger.Fatalf("gRPC server failed to serve: %v", err)
	}
}

func (gs *grpcServer) GetManifest(
	ctx context.Context,
	in *pb.ManifestReq,
) (*manifest.Manifest, error) {
	return gs.service.manifest, nil
}

func (gs *grpcServer) DoRequest(
	ctx context.Context,
	in *pb.Request,
) (*pb.Response, error) {
	req := &Request{
		HttpMethod:  in.HttpMethod.String(),
		Handler:     in.Handler,
		Path:        in.Path,
		Headers:     decodeMapStringArray(in.Headers),
		PathParams:  in.PathParams,
		QueryParams: decodeMapStringArray(in.QueryParams),
	}
	if req.HttpMethod == "PUT" || req.HttpMethod == "POST" {
		var body map[string]interface{}
		err := json.Unmarshal(in.Body, &body)
		if err != nil {
			return &pb.Response{
				Status: 400,
				Body:   []byte(`{"title":"Invalid JSON specified in body"}`),
			}, nil
		}
		req.Body = body
	}
	res := &Response{Body: make(map[string]interface{})}
	gs.service.doRequest(req, res)
	resBody, err := json.Marshal(res.Body)
	if err != nil {
		// There is a bug in the service code, panicking is the best course of action here
		// so that this (hopefully) never makes it into production.
		gs.logger.Panicf("failed to marshal json: %v", err)
	}
	return &pb.Response{
		Status:  int32(res.StatusCode),
		Body:    resBody,
		Headers: encodeMapStringArray(res.Headers),
	}, nil
}

func encodeMapStringArray(m map[string][]string) map[string]*pb.StringArray {
	result := map[string]*pb.StringArray{}
	for key, value := range m {
		result[key] = &pb.StringArray{Values: value}
	}
	return result
}

func decodeMapStringArray(m map[string]*pb.StringArray) map[string][]string {
	result := map[string][]string{}
	for key, value := range m {
		result[key] = value.Values
	}
	return result
}

/*
TODO
How to AWS implement the type conversion for lambbda handlers? That' is quite nice.
It would be good to do the same thing for handler.
*/
