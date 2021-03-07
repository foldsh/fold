import { Map as ProtoMap } from "google-protobuf";

import { HTTPMethod } from "../http";
import { Manifest, Route, Version } from "../../dist/proto/manifest_pb";
import {
  FoldHTTPMethod,
  FoldHTTPMethodMap,
  StringArray,
} from "../../dist/proto/http_pb";
import { Version as ServiceVersion } from "../version";

export const CHARSET_REGEXP = /;\s*charset\s*=/;
export const HEADER_SEPARATOR = ", ";
export const COOKIE_SEPARATOR = "; ";

export function splitHeader(header: string, value: string): string[] {
  const separator = isSetCookie(header) ? COOKIE_SEPARATOR : HEADER_SEPARATOR;
  return value.split(separator);
}

export function joinHeader(header: string, values: string[]) {
  const separator = isSetCookie(header) ? COOKIE_SEPARATOR : HEADER_SEPARATOR;
  return values.join(separator);
}

function isSetCookie(header: string): boolean {
  return header.toLowerCase() === "set-cookie";
}

export function decodeFoldHTTPMethod(
  n: FoldHTTPMethodMap[keyof FoldHTTPMethodMap]
): HTTPMethod {
  switch (n) {
    case FoldHTTPMethod.GET:
      return HTTPMethod.GET;

    case FoldHTTPMethod.HEAD:
      return HTTPMethod.HEAD;

    case FoldHTTPMethod.POST:
      return HTTPMethod.POST;

    case FoldHTTPMethod.PUT:
      return HTTPMethod.PUT;

    case FoldHTTPMethod.DELETE:
      return HTTPMethod.DELETE;

    case FoldHTTPMethod.CONNECT:
      return HTTPMethod.CONNECT;

    case FoldHTTPMethod.OPTIONS:
      return HTTPMethod.OPTIONS;

    case FoldHTTPMethod.TRACE:
      return HTTPMethod.TRACE;

    case FoldHTTPMethod.PATCH:
      return HTTPMethod.PATCH;

    default:
      // This can't happen actually happen in practice as it has all been
      // validated by the gRPC server.
      throw new Error(`Invalid HTTP method ${n}`);
  }
}

export function decodeHeaders(
  map: ProtoMap<string, StringArray>
): { [key: string]: string | string[] } {
  let results: { [key: string]: string | string[] } = {};
  map.forEach((value: StringArray, key: string) => {
    // Node standardises on lower case headers
    const name = key.toLowerCase();
    const values = value.getValuesList();
    // See https://nodejs.org/api/http.html#http_message_headers
    // Everything except set-cookie is a string joined with commas
    if (name === "set-cookie") {
      results[name] = values;
    } else {
      results[name] = values.join(HEADER_SEPARATOR);
    }
  });
  return results;
}

export function decodeQueryParams(
  map: ProtoMap<string, StringArray>
): { [key: string]: string | string[] } {
  let results: { [key: string]: string | string[] } = {};
  map.forEach((value: StringArray, key: string) => {
    const values = value.getValuesList();
    results[key] = values.length === 1 ? values[0] : values;
  });
  return results;
}

export function decodePathParams(
  map: ProtoMap<string, string>
): { [key: string]: string } {
  let results: { [key: string]: string } = {};
  map.forEach((value: string, key: string) => {
    results[key] = value;
  });
  return results;
}

export interface ManifestSpec {
  name?: string;
  version?: ServiceVersion;
  routes?: RouteSpec[];
}

export interface RouteSpec {
  method: string;
  handler: string;
  route: string;
}

export function makeManifest(manifest: ManifestSpec): Manifest {
  const params = newManifest();
  if (manifest.name) {
    params.name = manifest.name;
  }
  if (manifest.version) {
    params.version = manifest.version;
  }
  if (manifest.routes) {
    for (const route of manifest.routes) {
      params.routesList.push({
        httpMethod: FoldHTTPMethod[route.method as keyof FoldHTTPMethodMap],
        route: route.route,
      });
    }
  }
  return manifestFromObject(params);
}

function newManifest(): Manifest.AsObject {
  return ({ routesList: [] } as unknown) as Manifest.AsObject;
}

function manifestFromObject(obj: Manifest.AsObject): Manifest {
  const m = new Manifest();
  m.setName(obj.name);
  if (obj.version) {
    m.setVersion(versionFromObject(obj.version));
  }
  if (obj.routesList) {
    obj.routesList.forEach((route) => {
      m.addRoutes(routeFromObject(route));
    });
  }
  return m;
}

function versionFromObject(obj: Version.AsObject): Version {
  const r = new Version();
  r.setMajor(obj.major);
  r.setMinor(obj.minor);
  r.setPatch(obj.patch);
  return r;
}

function routeFromObject(obj: Route.AsObject): Route {
  const r = new Route();
  r.setHttpMethod(obj.httpMethod);
  r.setRoute(obj.route);
  return r;
}
