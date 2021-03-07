// package: http
// file: http.proto

import * as jspb from "google-protobuf";

export class FoldHTTPRequest extends jspb.Message {
  getHttpMethod(): FoldHTTPMethodMap[keyof FoldHTTPMethodMap];
  setHttpMethod(value: FoldHTTPMethodMap[keyof FoldHTTPMethodMap]): void;

  getPath(): string;
  setPath(value: string): void;

  getRawQuery(): string;
  setRawQuery(value: string): void;

  getFragment(): string;
  setFragment(value: string): void;

  hasHttpProto(): boolean;
  clearHttpProto(): void;
  getHttpProto(): FoldHTTPProto | undefined;
  setHttpProto(value?: FoldHTTPProto): void;

  getHost(): string;
  setHost(value: string): void;

  getRemoteAddr(): string;
  setRemoteAddr(value: string): void;

  getRequestUri(): string;
  setRequestUri(value: string): void;

  getContentLength(): number;
  setContentLength(value: number): void;

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
  getRoute(): string;
  setRoute(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FoldHTTPRequest.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FoldHTTPRequest
  ): FoldHTTPRequest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FoldHTTPRequest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FoldHTTPRequest;
  static deserializeBinaryFromReader(
    message: FoldHTTPRequest,
    reader: jspb.BinaryReader
  ): FoldHTTPRequest;
}

export namespace FoldHTTPRequest {
  export type AsObject = {
    httpMethod: FoldHTTPMethodMap[keyof FoldHTTPMethodMap];
    path: string;
    rawQuery: string;
    fragment: string;
    httpProto?: FoldHTTPProto.AsObject;
    host: string;
    remoteAddr: string;
    requestUri: string;
    contentLength: number;
    body: Uint8Array | string;
    headersMap: Array<[string, StringArray.AsObject]>;
    pathParamsMap: Array<[string, string]>;
    queryParamsMap: Array<[string, StringArray.AsObject]>;
    route: string;
  };
}

export class FoldHTTPResponse extends jspb.Message {
  getStatus(): number;
  setStatus(value: number): void;

  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): void;

  getHeadersMap(): jspb.Map<string, StringArray>;
  clearHeadersMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FoldHTTPResponse.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FoldHTTPResponse
  ): FoldHTTPResponse.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FoldHTTPResponse,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FoldHTTPResponse;
  static deserializeBinaryFromReader(
    message: FoldHTTPResponse,
    reader: jspb.BinaryReader
  ): FoldHTTPResponse;
}

export namespace FoldHTTPResponse {
  export type AsObject = {
    status: number;
    body: Uint8Array | string;
    headersMap: Array<[string, StringArray.AsObject]>;
  };
}

export class FoldHTTPProto extends jspb.Message {
  getProto(): string;
  setProto(value: string): void;

  getMajor(): number;
  setMajor(value: number): void;

  getMinor(): number;
  setMinor(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FoldHTTPProto.AsObject;
  static toObject(
    includeInstance: boolean,
    msg: FoldHTTPProto
  ): FoldHTTPProto.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: FoldHTTPProto,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): FoldHTTPProto;
  static deserializeBinaryFromReader(
    message: FoldHTTPProto,
    reader: jspb.BinaryReader
  ): FoldHTTPProto;
}

export namespace FoldHTTPProto {
  export type AsObject = {
    proto: string;
    major: number;
    minor: number;
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

export interface FoldHTTPMethodMap {
  GET: 0;
  HEAD: 1;
  POST: 2;
  PUT: 3;
  DELETE: 4;
  CONNECT: 5;
  OPTIONS: 6;
  TRACE: 7;
  PATCH: 8;
}

export const FoldHTTPMethod: FoldHTTPMethodMap;
