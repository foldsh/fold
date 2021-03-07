import { IncomingMessage } from "http";
import { Socket, isIP } from "net";
import { ParsedUrlQuery } from "querystring";

import { Ranges } from "range-parser";
import parseRange from "range-parser";
import accepts from "accepts";
import typeIs from "type-is";
import fresh from "fresh";

import { FoldHTTPRequest } from "../../dist/proto/http_pb";
import {
  HEADER_SEPARATOR,
  decodeFoldHTTPMethod,
  decodeQueryParams,
  decodePathParams,
  decodeHeaders,
} from "./utils";

import { Request, AcceptsReturns, Protocol, Response } from "../http";

/**
 * The job of FoldHTTPRequestWrapper is to implement the http.Request
 * interface for the FoldHTTPRequest.
 */
export class FoldHTTPRequestWrapper extends IncomingMessage implements Request {
  private readonly req: FoldHTTPRequest;
  private readonly _params: { [key: string]: string };
  private readonly _query: { [key: string]: string | string[] };
  private readonly _originalUrl: string;

  private _res!: Response;
  private _cookies: any;
  private _signedCookies: any;
  private _secret: string | undefined;
  private _path: string;
  private _timedout: boolean;

  constructor(req: FoldHTTPRequest) {
    super(new Socket({ readable: true }));
    // Push the body into the socket so it can be read with the typical NodeJS
    // stream interface.
    this.push(req.getBody_asU8(), "utf-8");
    // Pushing null indicates that the stream has ended and sets up the socket
    // to emit an 'end' event once the main body has been read.
    this.push(null);
    this.req = req;
    this.headers = decodeHeaders(this.req.getHeadersMap());
    this._query = decodeQueryParams(this.req.getQueryParamsMap());
    this._params = decodePathParams(this.req.getPathParamsMap());
    this.method = decodeFoldHTTPMethod(this.req.getHttpMethod());
    const proto = this.req.getHttpProto();
    if (proto) {
      this.httpVersion = proto.getProto();
      this.httpVersionMajor = proto.getMajor();
      this.httpVersionMinor = proto.getMinor();
    }
    this._cookies = undefined;
    this._signedCookies = undefined;
    this._secret = undefined;
    this._path = this.req.getPath();
    const rawQuery = this.req.getRawQuery();
    const url = `${this.path}${rawQuery ? "?" + rawQuery : ""}`;
    // The purpose of orginalUrl is to preserve the original url without any rewrites.
    this._originalUrl = url;
    this.url = url;
    this._timedout = false;
  }

  get res(): Response {
    return this._res;
  }

  set res(res: Response) {
    this._res = res;
  }

  accepts(...types: string[]): AcceptsReturns {
    const accept = accepts(this);
    const result = accept.type(types);
    if (Array.isArray(result)) {
      return result[0];
    }
    return result;
  }

  acceptsCharsets(...charsets: string[]): AcceptsReturns {
    const accept = accepts(this);
    return accept.charsets(charsets);
  }

  acceptsEncodings(...encodings: string[]): AcceptsReturns {
    const accept = accepts(this);
    return accept.encodings(encodings);
  }

  acceptsLanguages(...languages: string[]): AcceptsReturns {
    const accept = accepts(this);
    return accept.languages(languages);
  }

  get(header: string): string | string[] | undefined {
    return this.headers[header.toLowerCase()];
  }

  header(header: string): string | string[] | undefined {
    return this.headers[header];
  }

  get hostname(): string | undefined {
    // The runtime will always define the Host on the request.
    return this.req.getHost();
  }

  is(...types: string[]): string | false | null {
    return typeIs(this, types);
  }

  get originalUrl(): string {
    return this._originalUrl;
  }

  get path(): string {
    return this._path;
  }

  set path(path: string) {
    this._path = path;
  }

  get protocol(): Protocol {
    // The runtime will always define the url.
    return "https";
  }

  get query(): ParsedUrlQuery {
    return this._query;
  }
  get params(): { [key: string]: string } {
    return this._params;
  }

  get route(): string {
    return this.req.getRoute();
  }

  set route(route: string) {
    this.req.setRoute(route);
  }

  range(size: number, options?: any): -1 | -2 | Ranges | undefined {
    const range = this.get("Range") as string;
    if (!range) {
      return;
    }
    return parseRange(size, range, options);
  }

  get secure(): boolean {
    return this.protocol === "https";
  }

  get xhr(): boolean {
    const xRequestedWith = this.get("x-requested-with") as string | undefined;
    return xRequestedWith === undefined
      ? false
      : xRequestedWith === "XMLHttpRequest";
  }

  get ip(): string {
    let xForwardedFor = this.xForwardedFor;
    if (xForwardedFor.length === 0) {
      return "";
    }
    return xForwardedFor[0];
  }
  get ips(): string[] {
    return this.xForwardedFor;
  }

  private get xForwardedFor(): string[] {
    const xForwardedFor = this.get("x-forwarded-for") as string | undefined;
    if (xForwardedFor === undefined) {
      return [];
    }
    return xForwardedFor.split(HEADER_SEPARATOR);
  }

  get subdomains(): string[] {
    const hostname = this.hostname;

    if (!hostname) return [];

    // var offset = this.app.get("subdomain offset");
    // TODO need to app level config
    const offset = 2;
    const subdomains = !isIP(hostname)
      ? hostname.split(".").reverse()
      : [hostname];

    return subdomains.slice(offset);
  }
  get cookies(): any {
    return this._cookies;
  }

  set cookies(cookies: any) {
    this._cookies = cookies;
  }

  get signedCookies(): any {
    return this._signedCookies;
  }

  set signedCookies(signedCookies: any) {
    this._signedCookies = signedCookies;
  }

  get secret(): any {
    return this._secret;
  }

  set secret(secret: any) {
    this._secret = secret;
  }

  get fresh(): boolean {
    const method = this.method;
    const res = this.res;
    const status = res.statusCode;

    // GET or HEAD for weak freshness validation only
    if ("GET" !== method && "HEAD" !== method) return false;

    // 2xx or 304 as per rfc2616 14.26
    if ((status >= 200 && status < 300) || 304 === status) {
      return fresh(this.headers, {
        etag: res.get("ETag"),
        "last-modified": res.get("Last-Modified"),
      });
    }
    return false;
  }

  get stale(): boolean {
    return !this.fresh;
  }

  get timedout(): boolean {
    return this._timedout;
  }

  set timedout(timedout: boolean) {
    this._timedout = timedout;
  }
}
