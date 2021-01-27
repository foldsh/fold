// package: ingress
// file: ingress.proto

import * as jspb from "google-protobuf";
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

export class Request extends jspb.Message {
  getHttpMethod(): manifest_pb.HttpMethodMap[keyof manifest_pb.HttpMethodMap];
  setHttpMethod(
    value: manifest_pb.HttpMethodMap[keyof manifest_pb.HttpMethodMap]
  ): void;

  getHandler(): string;
  setHandler(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): void;

  getHeadersMap(): jspb.Map<string, StringArray>;
  clearHeadersMap(): void;
  getPathParamsMap(): jspb.Map<string, string>;
  clearPathParamsMap(): void;
  getQueryParamsMap(): jspb.Map<string, StringArray>;
  clearQueryParamsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Request.AsObject;
  static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Request,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Request;
  static deserializeBinaryFromReader(
    message: Request,
    reader: jspb.BinaryReader
  ): Request;
}

export namespace Request {
  export type AsObject = {
    httpMethod: manifest_pb.HttpMethodMap[keyof manifest_pb.HttpMethodMap];
    handler: string;
    path: string;
    body: Uint8Array | string;
    headersMap: Array<[string, StringArray.AsObject]>;
    pathParamsMap: Array<[string, string]>;
    queryParamsMap: Array<[string, StringArray.AsObject]>;
  };
}

export class StringArray extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<string>;
  setValuesList(value: Array<string>): void;
  addValues(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringArray.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: StringArray
  ): StringArray.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: StringArray,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): StringArray;
  static deserializeBinaryFromReader(
    message: StringArray,
    reader: jspb.BinaryReader
  ): StringArray;
}

export namespace StringArray {
  export type AsObject = {
    valuesList: Array<string>;
  };
}

export class Response extends jspb.Message {
  getStatus(): number;
  setStatus(value: number): void;

  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): void;

  getHeadersMap(): jspb.Map<string, StringArray>;
  clearHeadersMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Response.AsObject;
  static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Response,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Response;
  static deserializeBinaryFromReader(
    message: Response,
    reader: jspb.BinaryReader
  ): Response;
}

export namespace Response {
  export type AsObject = {
    status: number;
    body: Uint8Array | string;
    headersMap: Array<[string, StringArray.AsObject]>;
  };
}
