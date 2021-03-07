// package: manifest
// file: manifest.proto

import * as jspb from "google-protobuf";
import * as http_pb from "./http_pb";

export class Manifest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  hasVersion(): boolean;
  clearVersion(): void;
  getVersion(): Version | undefined;
  setVersion(value?: Version): void;

  hasBuildInfo(): boolean;
  clearBuildInfo(): void;
  getBuildInfo(): BuildInfo | undefined;
  setBuildInfo(value?: BuildInfo): void;

  clearRoutesList(): void;
  getRoutesList(): Array<Route>;
  setRoutesList(value: Array<Route>): void;
  addRoutes(value?: Route, index?: number): Route;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Manifest.AsObject;
  static toObject(includeInstance: boolean, msg: Manifest): Manifest.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Manifest,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Manifest;
  static deserializeBinaryFromReader(
    message: Manifest,
    reader: jspb.BinaryReader
  ): Manifest;
}

export namespace Manifest {
  export type AsObject = {
    name: string;
    version?: Version.AsObject;
    buildInfo?: BuildInfo.AsObject;
    routesList: Array<Route.AsObject>;
  };
}

export class BuildInfo extends jspb.Message {
  getMaintainer(): string;
  setMaintainer(value: string): void;

  getImage(): string;
  setImage(value: string): void;

  getTag(): string;
  setTag(value: string): void;

  getPath(): string;
  setPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BuildInfo.AsObject;
  static toObject(includeInstance: boolean, msg: BuildInfo): BuildInfo.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: BuildInfo,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): BuildInfo;
  static deserializeBinaryFromReader(
    message: BuildInfo,
    reader: jspb.BinaryReader
  ): BuildInfo;
}

export namespace BuildInfo {
  export type AsObject = {
    maintainer: string;
    image: string;
    tag: string;
    path: string;
  };
}

export class Version extends jspb.Message {
  getMajor(): number;
  setMajor(value: number): void;

  getMinor(): number;
  setMinor(value: number): void;

  getPatch(): number;
  setPatch(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Version.AsObject;
  static toObject(includeInstance: boolean, msg: Version): Version.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Version,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Version;
  static deserializeBinaryFromReader(
    message: Version,
    reader: jspb.BinaryReader
  ): Version;
}

export namespace Version {
  export type AsObject = {
    major: number;
    minor: number;
    patch: number;
  };
}

export class Route extends jspb.Message {
  getHttpMethod(): http_pb.FoldHTTPMethodMap[keyof http_pb.FoldHTTPMethodMap];
  setHttpMethod(
    value: http_pb.FoldHTTPMethodMap[keyof http_pb.FoldHTTPMethodMap]
  ): void;

  getRoute(): string;
  setRoute(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Route.AsObject;
  static toObject(includeInstance: boolean, msg: Route): Route.AsObject;
  static extensions: { [key: number]: jspb.ExtensionFieldInfo<jspb.Message> };
  static extensionsBinary: {
    [key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>;
  };
  static serializeBinaryToWriter(
    message: Route,
    writer: jspb.BinaryWriter
  ): void;
  static deserializeBinary(bytes: Uint8Array): Route;
  static deserializeBinaryFromReader(
    message: Route,
    reader: jspb.BinaryReader
  ): Route;
}

export namespace Route {
  export type AsObject = {
    httpMethod: http_pb.FoldHTTPMethodMap[keyof http_pb.FoldHTTPMethodMap];
    route: string;
  };
}
