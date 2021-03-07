// GENERATED CODE -- DO NOT EDIT!

// Original file comments:
// This defines the interface between the fold runtime and an application
// for inbound traffic. I.e., it manages inbound requests or events on the
// way in and passes them on to application.
"use strict";
var grpc = require("@grpc/grpc-js");
var ingress_pb = require("./ingress_pb.js");
var http_pb = require("./http_pb.js");
var manifest_pb = require("./manifest_pb.js");

function serialize_http_FoldHTTPRequest(arg) {
  if (!(arg instanceof http_pb.FoldHTTPRequest)) {
    throw new Error("Expected argument of type http.FoldHTTPRequest");
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_http_FoldHTTPRequest(buffer_arg) {
  return http_pb.FoldHTTPRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_http_FoldHTTPResponse(arg) {
  if (!(arg instanceof http_pb.FoldHTTPResponse)) {
    throw new Error("Expected argument of type http.FoldHTTPResponse");
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_http_FoldHTTPResponse(buffer_arg) {
  return http_pb.FoldHTTPResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_ingress_ManifestReq(arg) {
  if (!(arg instanceof ingress_pb.ManifestReq)) {
    throw new Error("Expected argument of type ingress.ManifestReq");
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_ingress_ManifestReq(buffer_arg) {
  return ingress_pb.ManifestReq.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_manifest_Manifest(arg) {
  if (!(arg instanceof manifest_pb.Manifest)) {
    throw new Error("Expected argument of type manifest.Manifest");
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_manifest_Manifest(buffer_arg) {
  return manifest_pb.Manifest.deserializeBinary(new Uint8Array(buffer_arg));
}

var FoldIngressService = (exports.FoldIngressService = {
  // Retrieve the manifest from the service.
  getManifest: {
    path: "/ingress.FoldIngress/GetManifest",
    requestStream: false,
    responseStream: false,
    requestType: ingress_pb.ManifestReq,
    responseType: manifest_pb.Manifest,
    requestSerialize: serialize_ingress_ManifestReq,
    requestDeserialize: deserialize_ingress_ManifestReq,
    responseSerialize: serialize_manifest_Manifest,
    responseDeserialize: deserialize_manifest_Manifest,
  },
  // Ask the service to process an HTTP request.
  doRequest: {
    path: "/ingress.FoldIngress/DoRequest",
    requestStream: false,
    responseStream: false,
    requestType: http_pb.FoldHTTPRequest,
    responseType: http_pb.FoldHTTPResponse,
    requestSerialize: serialize_http_FoldHTTPRequest,
    requestDeserialize: deserialize_http_FoldHTTPRequest,
    responseSerialize: serialize_http_FoldHTTPResponse,
    responseDeserialize: deserialize_http_FoldHTTPResponse,
  },
});

exports.FoldIngressClient = grpc.makeGenericClientConstructor(
  FoldIngressService
);
