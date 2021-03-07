import { OutgoingHttpHeader, OutgoingHttpHeaders, ServerResponse } from "http";

import contentType from "content-type";
import mime from "mime";
import encodeUrl from "encodeurl";
import cookie from "cookie";
import { sign } from "cookie-signature";
import merge from "utils-merge";
import send from "send";
import status from "statuses";
import vary from "vary";

import { Response, SerializeOptions } from "../http";
import { FoldHTTPResponse, StringArray } from "../../dist/proto/http_pb";
import { FoldHTTPRequestWrapper } from "./req-wrapper";
import {
  CHARSET_REGEXP,
  HEADER_SEPARATOR,
  joinHeader,
  splitHeader,
} from "./utils";

/**
 * The job of FoldHTTPResponseWrapper is to implement the http.Response
 * interface for the FoldHTTPResponse.
 */
export class FoldHTTPResponseWrapper
  extends ServerResponse
  implements Response {
  req: FoldHTTPRequestWrapper;
  private res: FoldHTTPResponse;

  constructor(req: FoldHTTPRequestWrapper, res: FoldHTTPResponse) {
    super(req);
    this.socket = req.socket;
    this.req = req;
    this.res = res;
  }

  append(field: string, value: any): Response {
    const prev = this.get(field);
    let newValue = value;

    if (prev) {
      if (Array.isArray(prev)) {
        newValue = prev.concat(newValue);
      } else if (Array.isArray(newValue)) {
        newValue = joinHeader(field, [prev].concat(newValue) as string[]);
      } else {
        newValue = joinHeader(field, [prev, newValue]);
      }
    } else if (Array.isArray(newValue)) {
      newValue = joinHeader(field, newValue);
    }

    return this.set(field, newValue);
  }

  /**
   * Clears a cookie by name. The options must be the same as the options
   * used to set the cookie in the first place.
   */
  clearCookie(name: string, options?: Partial<SerializeOptions>): Response {
    let opts = merge(
      { expires: new Date(1), path: "/" },
      options
    ) as SerializeOptions;
    return this.cookie(name, "", opts);
  }

  /**
   * Set cookie `name` to `value`, with the given `options`.
   *
   * Options:
   *
   *    - `maxAge`   max-age in milliseconds, converted to `expires`
   *    - `signed`   sign the cookie
   *    - `path`     defaults to "/"
   *    - `encode`   a function to encode a cookie's value
   *    - `maxAge`   sets the Max-Age attribute (in seconds)
   *    - `domain`   set the Domain attribute
   *    - `httpOnly` set the HttpOnly attribute
   *    - `secure`   set the Secure attribute
   *    - `sameSite` set the SameSite attribute
   *                 true: Strict
   *                 false: don't set Samesite
   *                 "lax": Lax
   *                 "strict": Strict
   *                 "none": None
   *    - `expires`  expirty date for the cookie
   *
   * Examples:
   *
   *    // "Remember Me" for 15 minutes
   *    res.cookie('rememberme', '1', { expires: new Date(Date.now() + 900000), httpOnly: true });
   *
   *    // same as above
   *    res.cookie('rememberme', '1', { maxAge: 900000, httpOnly: true })
   */
  cookie(
    name: string,
    value: string | Record<string, unknown>,
    options?: Partial<SerializeOptions> & Partial<{ signed: boolean }>
  ): Response {
    options = merge(options ? options : {}, {}) as SerializeOptions;
    let secret = this.req.secret;
    let signed = options.signed;

    if (signed && !secret) {
      throw new Error('cookieParser("secret") required for signed cookies');
    }

    let val =
      typeof value === "object" ? "j:" + JSON.stringify(value) : String(value);

    if (signed) {
      val = "s:" + sign(val, secret);
    }

    if (options.maxAge) {
      options.expires = new Date(Date.now() + options.maxAge);
      options.maxAge /= 1000;
    }

    if (options.path == null) {
      options.path = "/";
    }

    return this.append(
      "Set-Cookie",
      cookie.serialize(name, String(val), options)
    );
  }

  get(field: string): string | number | string[] | undefined {
    return this.getHeader(field);
  }

  header(
    field: string | Record<string, string | string[]>,
    val?: string | string[]
  ): Response {
    return this.set(field, val);
  }

  set(
    field: string | Record<string, string | string[]>,
    val?: string | string[]
  ): Response {
    if (arguments.length === 2) {
      // Handle the set(header, value) case.
      if (typeof field !== "string") {
        throw new TypeError("Field must be of type string");
      }
      let value: string | string[];
      const isArray = Array.isArray(val);
      // Node standardises to lower case headers
      field = field.toLowerCase();
      // // Only set-cookie can be an array.
      if (isArray && field !== "set-cookie") {
        throw new TypeError("Only Set-Cookie can be set to an Array");
      } else if (isArray) {
        value = (val as string[]).map(String);
      } else {
        value = String(val);
      }
      // add charset to content-type
      if (field === "content-type") {
        if (Array.isArray(value)) {
          throw new TypeError("Content-Type cannot be set to an Array");
        }
        if (!CHARSET_REGEXP.test(value)) {
          var charset = send.mime.charsets.lookup(value.split(";")[0], "");
          if (charset) value += "; charset=" + charset.toLowerCase();
        }
      }
      this.setHeader(field as string, value);
    } else {
      // Handle the {header: value} case.
      for (let [header, value] of Object.entries(field)) {
        this.set(header, value);
      }
    }
    return this;
  }

  /**
   * Set Link header field with the given `links`.
   *
   * Examples:
   *
   *    res.links({
   *      next: 'http://api.example.com/users?page=2',
   *      last: 'http://api.example.com/users?page=5'
   *    });
   *
   */
  links(links: { [p: string]: string }): Response {
    let link = this.get("Link") as string;
    const newLinks: string[] = [];
    Object.keys(links).forEach(function (rel) {
      newLinks.push(`<${links[rel]}>; rel="${rel}"`);
    });
    const newLinkString = joinHeader("Link", newLinks);
    if (link) {
      return this.set("Link", `${link}${HEADER_SEPARATOR}${newLinkString}`);
    } else {
      return this.set("Link", newLinkString);
    }
  }

  location(url: string): Response {
    var loc = url;

    // "back" is an alias for the referrer
    if (url === "back") {
      let refererr = this.req.get("Referrer");
      if (!refererr) {
        loc = "/";
      } else {
        loc = refererr as string;
      }
    }

    // set location
    return this.set("Location", encodeUrl(loc));
  }

  redirect(url: string, status?: number): Response {
    if (!status) {
      this.status(302);
    } else {
      this.status(status);
    }
    this.location(url);
    return this.json({ title: `Redirecting to ${url}` });
  }

  json(body: unknown): Response {
    this.type("application/json");
    this.send(JSON.stringify(body));
    return this;
  }

  /**
   * Our job is fairly simple here, we just care about getting data into
   * the response. Content length, headers etc are handled by the runtime.
   * @param body The data you wish to write as the response body.
   */
  send(body: unknown): Response {
    // TODO etags
    // TODO freshness
    // TODO strip unnecessary headers
    if (body === null) {
      return this;
    }
    switch (typeof body) {
      case "string":
        if (!this.get("Content-Type")) {
          this.type("text");
        }
        let type = this.get("Content-Type") as string;
        if (type) {
          this.type(setCharset(type, "utf-8"));
        }
        this.end(Buffer.from(body));
        break;
      case "boolean":
      case "number":
      case "object":
        if (Buffer.isBuffer(body)) {
          if (!this.get("Content-Type")) {
            this.type("bin");
            this.end(body);
          }
        } else {
          return this.json(body);
        }
        break;
    }
    return this;
  }

  // @ts-ignore
  // No idea how to make TS happy with this, they are copied from the
  // stdlib type signatures.
  end(cb?: () => void): void;
  end(chunk: any, cb?: () => void): void;
  end(chunk: any, encoding: BufferEncoding, cb?: () => void): void {
    // Sort out the arguments.
    let _chunk = undefined;
    let _encoding = undefined;
    let _callback = undefined;

    if (arguments.length == 1) {
      // The first signature
      if (typeof chunk === "function") {
        _callback = chunk;
      } else {
        _chunk = chunk;
      }
    } else if (arguments.length == 2) {
      // Either the second signature or the third with no callback.
      _chunk = chunk;
      if (typeof encoding == "function") {
        // It's the third signature with no cb set.
        _callback = encoding;
      } else {
        _encoding = _encoding;
        _callback = cb;
      }
    } else {
      // All three params are set.
      _chunk = chunk;
      _encoding = encoding;
      _callback = cb;
    }
    if (!_encoding) {
      _encoding = "utf-8";
    }

    super.end(_chunk, _encoding as BufferEncoding, _callback);

    if (_chunk) {
      // Rather than mess around with the socket, we are just adding
      // the body directly.
      this.res.setBody(chunk);
    }
    this.emit("finish");
  }

  writeHead(
    statusCode: number,
    reasonPhrase?: string | (OutgoingHttpHeaders | OutgoingHttpHeader[]),
    headers?: OutgoingHttpHeaders | OutgoingHttpHeader[]
  ): this {
    // Sort out the arguments.
    let _reasonPhrase: string | undefined = undefined;
    let _headers: (OutgoingHttpHeaders | OutgoingHttpHeader[]) | undefined;
    if (arguments.length == 2) {
      if (typeof reasonPhrase === "string") {
        _reasonPhrase = reasonPhrase;
      } else {
        _headers = reasonPhrase;
      }
    } else {
      _reasonPhrase = _reasonPhrase;
      _headers = headers;
    }

    // Now we add all the headers in to our map and set the status code.
    this.statusCode = statusCode;
    if (_headers) {
      if (Array.isArray(_headers)) {
        throw new TypeError(
          "Fold does not support writing headers in as a string."
        );
      }
      this.set(_headers as Record<string, string | string[]>);
    }
    // And finally return with the call to super.
    // It' useful to do this so that the state of the Request object
    // matches the state of the fold response object we're building.
    // This makes it easier to work with third party middlewares as
    // the request behaves more like how they expect.
    return super.writeHead(statusCode, _reasonPhrase, _headers);
  }

  sendStatus(statusCode: number): Response {
    var body = status(statusCode) || String(statusCode);

    this.statusCode = statusCode;
    // this.res.setStatus(statusCode);
    this.type("txt");

    return this.send(body);
  }

  status(status: number): Response {
    this.statusCode = status;
    this.res.setStatus(status);
    return this;
  }

  type(type: string): Response {
    var ct = type.indexOf("/") === -1 ? mime.getType(type) : type;
    if (ct !== null) {
      this.set("Content-Type", ct);
    }
    return this;
  }

  vary(field: string): Response {
    vary(this, field);
    return this;
  }

  attachment(_filename?: string): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }
  sendFile(_path: string, _options?: any, _cb?: (err?: any) => void): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }
  render(_file: string, _data?: Record<string, any>, _options?: any): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }
  format(_obj: any): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }
  download(
    _path: string,
    _filename: string,
    _options?: any,
    _cb?: (err?: any) => void
  ): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }
  jsonp(_obj: any): Response {
    throw new Error(
      "Not yet implemented, please file an issue if this feature matters to you."
    );
  }

  // For internal use
  get foldHTTPResponse(): FoldHTTPResponse {
    // Update the response with the headers we've collected.
    let resHeaders = this.res.getHeadersMap();
    for (let [header, values] of Object.entries(this.getHeaders())) {
      let valueArray: string[];
      if (values === undefined) {
        continue;
      } else if (typeof values === "string") {
        valueArray = splitHeader(header, values);
      } else if (typeof values === "number") {
        valueArray = [String(values)];
      } else {
        valueArray = values;
      }
      let sa = new StringArray();
      sa.setValuesList(valueArray);
      resHeaders.set(header, sa);
    }
    // Update it with the status. It's best to do it here so that
    // libraries that update statusCode directly will still have
    // the desired effect.
    this.res.setStatus(this.statusCode);
    // Finally, lets check that the body has the appropriate data type for
    // gRPC. Bodies written diretly by us will already have done this
    // but third party middlewares can write strings and so on so we need
    // to encode those properly.
    if (typeof this.res.getBody() === "string") {
      this.res.setBody(Buffer.from(this.res.getBody() as string, "utf-8"));
    }
    return this.res;
  }
}

function setCharset(type: string, charset: string): string {
  if (!type || !charset) {
    return type;
  }
  var parsed = contentType.parse(type);
  parsed.parameters.charset = charset;
  return contentType.format(parsed);
}
