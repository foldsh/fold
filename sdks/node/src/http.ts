import { IncomingMessage, ServerResponse } from "http";
import { ParsedUrlQuery } from "querystring";
import { Ranges } from "range-parser";
import { Service } from "./service";

export enum HTTPMethod {
  GET = "GET",
  HEAD = "HEAD",
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE",
  CONNECT = "CONNECT",
  OPTIONS = "OPTIONS",
  TRACE = "TRACE",
  PATCH = "PATCH",
}

export type Protocol = "http" | "https" | string;

export type AcceptsReturns = string | false;

// export interface Request extends UnstreamableIncomingMessage {
export interface Request extends IncomingMessage {
  originalUrl: string;
  path: string;
  query: ParsedUrlQuery;
  params: { [key: string]: string };
  protocol: Protocol;
  secure: boolean;
  xhr: boolean;
  hostname: string | undefined;
  ip?: string;
  ips?: string[];
  subdomains?: string[];
  cookies?: any;
  signedCookies?: any;
  secret?: string;
  fresh?: boolean;
  stale?: boolean;
  body?: any;
  accepted?: boolean;
  param?: any;
  host?: string;
  route?: string;
  app?: Service;
  baseUrl?: string;
  timedout?: boolean;

  get: (header: string) => string | string[] | undefined;
  header: (header: string) => string | string[] | undefined;
  range: (size: number, options?: any) => -1 | -2 | Ranges | undefined;
  accepts: (...types: string[]) => AcceptsReturns;
  acceptsEncodings: (...encodings: string[]) => AcceptsReturns;
  acceptsCharsets: (...charsets: string[]) => AcceptsReturns;
  acceptsLanguages: (...languages: string[]) => AcceptsReturns;
  // Checks if the incoming request has one of the requested mime types
  is: (...types: string[]) => string | false | null;
}

/**
 * This chooses not to implement the following from the express
 * interface for now:
 *   - attachment
 *   - sendFile
 *   - render
 *   - format
 *   - download
 *
 * The focus is on JSON based REST APIs and these features are simply
 * not required for most tasks in that space. Multi part data in general
 * doesn't play well with how this all works at the moment. Then again,
 * lambda didn't support it for ages so it's clearly not going to stop
 * a product succeeding.
 */
// export interface Response extends UnstreamableServerResponse {
export interface Response extends ServerResponse {
  get(field: string): string | number | string[] | undefined;
  set(
    field: string | Record<string, string | string[]>,
    val?: string | string[]
  ): Response;
  header(
    field: string | Record<string, string | string[]>,
    val?: string | string[]
  ): Response;
  send(body: unknown): Response;
  json(body: unknown): Response;
  status(status: number): Response;
  sendStatus(statusCode: number): Response;
  cookie(
    name: string,
    value: string | Record<string, unknown>,
    options?: Partial<SerializeOptions> & Partial<{ signed: boolean }>
  ): Response;
  clearCookie(name: string, options?: Partial<SerializeOptions>): Response;
  location(url: string): Response;
  links(links: { [key: string]: string }): Response;
  vary(field: string): Response;
  redirect(url: string, status?: number): Response;
  type(type: string): Response;
  locals?: Record<string, any>;
  append(field: string, value: any): Response;
  // The following are not currently implemented, please file an issue
  // if you'd really like to use them!
  attachment(filename?: string): Response;
  sendFile(path: string, options?: any, cb?: (err?: any) => void): Response;
  render(file: string, data?: Record<string, any>, options?: any): Response;
  format(obj: any): Response;
  download(
    path: string,
    filename: string,
    options?: any,
    cb?: (err?: any) => void
  ): Response;
  jsonp(obj: any): Response;
}

/**
 * Options for serializing cookies.
 */
export interface SerializeOptions {
  encode?: (str: string) => string;
  maxAge?: number;
  domain?: string;
  path?: string;
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: true | false | "lax" | "strict" | "none";
  expires?: Date;
}
