// GENERATED CODE -- DO NOT EDIT!

// package: ingress
// file: ingress.proto

import * as ingress_pb from "./ingress_pb";
import * as manifest_pb from "./manifest_pb";
import * as grpc from "@grpc/grpc-js";

interface IFoldIngressService
  extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  getManifest: grpc.MethodDefinition<
    ingress_pb.ManifestReq,
    manifest_pb.Manifest
  >;
  doRequest: grpc.MethodDefinition<ingress_pb.Request, ingress_pb.Response>;
}

export const FoldIngressService: IFoldIngressService;

export interface IFoldIngressServer extends grpc.UntypedServiceImplementation {
  getManifest: grpc.handleUnaryCall<
    ingress_pb.ManifestReq,
    manifest_pb.Manifest
  >;
  doRequest: grpc.handleUnaryCall<ingress_pb.Request, ingress_pb.Response>;
}

export class FoldIngressClient extends grpc.Client {
  constructor(
    address: string,
    credentials: grpc.ChannelCredentials,
    options?: object
  );
  getManifest(
    argument: ingress_pb.ManifestReq,
    callback: grpc.requestCallback<manifest_pb.Manifest>
  ): grpc.ClientUnaryCall;
  getManifest(
    argument: ingress_pb.ManifestReq,
    metadataOrOptions: grpc.Metadata | grpc.CallOptions | null,
    callback: grpc.requestCallback<manifest_pb.Manifest>
  ): grpc.ClientUnaryCall;
  getManifest(
    argument: ingress_pb.ManifestReq,
    metadata: grpc.Metadata | null,
    options: grpc.CallOptions | null,
    callback: grpc.requestCallback<manifest_pb.Manifest>
  ): grpc.ClientUnaryCall;
  doRequest(
    argument: ingress_pb.Request,
    callback: grpc.requestCallback<ingress_pb.Response>
  ): grpc.ClientUnaryCall;
  doRequest(
    argument: ingress_pb.Request,
    metadataOrOptions: grpc.Metadata | grpc.CallOptions | null,
    callback: grpc.requestCallback<ingress_pb.Response>
  ): grpc.ClientUnaryCall;
  doRequest(
    argument: ingress_pb.Request,
    metadata: grpc.Metadata | null,
    options: grpc.CallOptions | null,
    callback: grpc.requestCallback<ingress_pb.Response>
  ): grpc.ClientUnaryCall;
}
