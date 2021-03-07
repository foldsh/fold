// package: ingress
// file: ingress.proto

import * as jspb from "google-protobuf";
import * as http_pb from "./http_pb";
import * as manifest_pb from "./manifest_pb";

export class ManifestReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ManifestReq.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: ManifestReq
  ): ManifestReq.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: ManifestReq,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): ManifestReq;
  static deserializeBinaryFromReader(
    message: ManifestReq,
    reader: jspb.BinaryReader
  ): ManifestReq;
}

export namespace ManifestReq {
  export type AsObject = {};
}
