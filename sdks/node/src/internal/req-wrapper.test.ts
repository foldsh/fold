import merge from "utils-merge";

import { FoldHTTPRequestWrapper } from "./req-wrapper";
import { Request } from "../http";
import { FoldHTTPRequest, FoldHTTPResponse } from "../../dist/proto/http_pb";
import {
  foldReqFromObject,
  foldResFromObject,
  newReq,
  headerMap,
  queryMap,
  paramsMap,
  newRes,
  updateWithURL,
} from "../__tests__/proto-utils";
import { FoldHTTPResponseWrapper } from "./res-wrapper";
import { HEADER_SEPARATOR } from "./utils";

describe("FoldHTTPRequestWrapper", () => {
  test("should get headers", () => {
    const req = makeReq({
      headersMap: headerMap({
        "Content-Type": ["application/json"],
        Links: ["foo", "bar"],
      }),
    });
    expect(req.header("content-type")).toEqual("application/json");
    expect(req.header("links")).toEqual("foo, bar");
  });

  test("should read body", () => {
    // Slightly odd test. Basically, express sets the body to undefined by default
    // unless you have used a middleware to parse it. Middleware is outside of the
    // scope of this test, so I am instead demonstrating that you can drain the
    // socket and get the body out of it.
    const content = `{"hello": "world"}`;
    const req = makeReq({
      body: content,
    });
    expect(req.body).toBeUndefined();
    // req.body = JSON.parse(Buffer.from(req.socket.read()).toString("utf-8"));
    req.body = JSON.parse(Buffer.from(req.read()).toString("utf-8"));
    expect(req.body).toEqual({ hello: "world" });
  });

  test("should verify Accept header", () => {
    const req = makeReq({
      headersMap: headerMap({
        Accept: ["text/*", "application/json"],
      }),
    });
    expect(req.accepts("html")).toEqual("html");
    expect(req.accepts("json")).toEqual("json");
    expect(req.accepts("application/json")).toEqual("application/json");
    expect(req.accepts("image/png")).toBe(false);
  });

  test("should verify Accept-Charset header", () => {
    const req = makeReq({
      headersMap: headerMap({
        "Accept-Charset": ["utf-8"],
      }),
    });
    expect(req.acceptsCharsets("utf-8")).toBe("utf-8");
    expect(req.acceptsCharsets("iso-8859-1")).toBe(false);
  });

  test("should verify Accept-Encoding header", () => {
    const req = makeReq({
      headersMap: headerMap({
        "Accept-Encoding": ["gzip", "deflate"],
      }),
    });
    expect(req.acceptsEncodings("gzip")).toEqual("gzip");
    expect(req.acceptsEncodings("deflate")).toEqual("deflate");
    // identity is always set, unless explicitly forbidden
    expect(req.acceptsEncodings("identity")).toEqual("identity");
    expect(req.acceptsEncodings("br")).toEqual(false);
  });

  test("should verify Accept-Language headers", () => {
    const req = makeReq({
      headersMap: headerMap({
        "Accept-Language": ["fr-CH", "fr;q=0.9", "en;q=0.8", "*;q=0.5"],
      }),
    });
    expect(req.acceptsLanguages("fr-CH")).toEqual("fr-CH");
    expect(req.acceptsLanguages("en", "fr")).toEqual("fr");
    expect(req.acceptsLanguages("de")).toEqual("de");
  });

  test("should get hostname", () => {
    const req = makeReq({ host: "www.fold.sh" });
    expect(req.hostname).toEqual("www.fold.sh");
  });

  test("should infer mimetype", () => {
    const req = makeReq({
      body: "<html>Hello, World!<html>",
      headersMap: headerMap({
        "Content-Type": ["text/html; charset=utf-8"],
      }),
    });
    expect(req.is("html")).toEqual("html");
    expect(req.is("text/html")).toEqual("text/html");
    expect(req.is("text/*")).toEqual("text/html");
    expect(req.is("json")).toBe(false);
  });

  test("should get range", () => {
    const req = makeReq({
      headersMap: headerMap({
        Range: ["bytes=0-50", "51-100"],
      }),
    });
    let range = req.range(21) as any[];
    expect(range.length).toEqual(1);
    expect(range[0]).toMatchObject({ start: 0, end: 20 });
  });

  test("should get url", () => {
    const req = makeReq({}, "https://www.fold.sh/test/suite?foo=bar");
    expect(req.url).toEqual("/test/suite?foo=bar");
  });
  test("should get orginal url", () => {
    const req = makeReq({}, "https://www.fold.sh/test/suite?foo=bar");
    expect(req.originalUrl).toEqual("/test/suite?foo=bar");
    // Rewriting url should preserve original url.
    req.url = "/suite";
    expect(req.originalUrl).toEqual("/test/suite?foo=bar");
  });

  test("should get path", () => {
    const req = makeReq({}, "https://www.fold.sh/test/suite?foo=bar");
    expect(req.path).toEqual("/test/suite");
  });

  test("should get protocol", () => {
    const req = makeReq({}, "https://www.fold.sh/test/suite?foo=bar");
    expect(req.protocol).toEqual("https");
  });

  test("should get http version", () => {
    const req = makeReq({
      httpProto: {
        proto: "HTTP/1.1",
        major: 1,
        minor: 1,
      },
    });
    expect(req.httpVersion).toEqual("HTTP/1.1");
    expect(req.httpVersionMajor).toEqual(1);
    expect(req.httpVersionMinor).toEqual(1);
  });

  test("should get method", () => {
    let req = makeReq({ httpMethod: 0 });
    expect(req.method).toEqual("GET");
    req = makeReq({ httpMethod: 1 });
    expect(req.method).toEqual("HEAD");
    req = makeReq({ httpMethod: 2 });
    expect(req.method).toEqual("POST");
  });

  test("should get secure", () => {
    let req = makeReq({});
    expect(req.secure).toBe(true);
  });

  test("should get query params", () => {
    let req = makeReq({});
    expect(req.query.names).toBe(undefined);
    req = makeReq({
      queryParamsMap: queryMap({ names: ["bill", "ted"], awesome: ["true"] }),
    });
    expect(req.query.names).toEqual(["bill", "ted"]);
    expect(req.query.awesome).toEqual("true");
  });

  test("should get path params", () => {
    let req = makeReq({});
    expect(req.params.id).toBe(undefined);
    req = makeReq({
      pathParamsMap: paramsMap({ id: "1234", name: "tom" }),
    });
    expect(req.params.id).toEqual("1234");
    expect(req.params.name).toEqual("tom");
  });

  test("should get ip params", () => {
    let req = makeReq({});
    expect(req.ip).toEqual("");
    req = makeReq({
      headersMap: headerMap({
        "X-Forwarded-For": ["203.0.113.195", "70.41.3.18", "150.172.238.178"],
      }),
    });
    expect(req.ip).toEqual("203.0.113.195");
  });

  test("should get ips params", () => {
    let req = makeReq({});
    expect(req.ips).toEqual([]);
    req = makeReq({
      headersMap: headerMap({
        "X-Forwarded-For": ["203.0.113.195", "70.41.3.18", "150.172.238.178"],
      }),
    });
    expect(req.ips).toEqual(["203.0.113.195", "70.41.3.18", "150.172.238.178"]);
  });

  test("should check if xhr", () => {
    let req = makeReq({});
    expect(req.xhr).toBe(false);
    req = makeReq({
      headersMap: headerMap({
        "X-Requested-With": ["XMLHttpRequest"],
      }),
    });
    expect(req.xhr).toBe(true);
    req = makeReq({
      headersMap: headerMap({
        "X-Requested-With": ["other-value"],
      }),
    });
    expect(req.xhr).toBe(false);
  });

  test("should get subdomains", () => {
    const req = makeReq({ host: "dev.mysubdomain.fold.sh" });
    expect(req.subdomains).toEqual(["mysubdomain", "dev"]);
  });

  test("should get cookies", () => {
    const req = makeReq({});
    // It is the job of middleware to set/parse cookies.
    // The request just needs to make the default undefined.
    expect(req.cookies).toEqual(undefined);
  });

  test("should get signedCookies", () => {
    const req = makeReq({});
    // It is the job of middleware to set/parse signedCookies.
    // The request just needs to make the default undefined.
    expect(req.signedCookies).toEqual(undefined);
  });

  test("should get secret", () => {
    const req = makeReq({});
    // It is the job of middleware to set/parse the secret.
    // The request just needs to make the default an empty object.
    expect(req.secret).toBe(undefined);
  });

  // Don't need to test all the cases as the fresh library handles it.
  // However I do want to make sure that I'm dealling with all the
  // parameters correctly so I've done a case for each of thoes.
  describe("req.fresh", () => {
    test("GET requests are not fresh", () => {
      const req = makeReq({ httpMethod: 0 }, "http://fold.sh", { status: 201 });
      expect(req.fresh).toBe(false);
    });
    test("HEAD requests are not fresh", () => {
      const req = makeReq({ httpMethod: 1 }, "http://fold.sh", { status: 201 });
      expect(req.fresh).toBe(false);
    });
    test(`requests with status 400 are not fresh`, () => {
      const req = makeReq({ httpMethod: 1 }, "http://fold.sh", {
        status: 400,
      });
      expect(req.fresh).toBe(false);
    });
    test(`requests with status 404 are not fresh`, () => {
      const req = makeReq({ httpMethod: 1 }, "http://fold.sh", {
        status: 404,
      });
      expect(req.fresh).toBe(false);
    });
    test(`requests with status 500 are not fresh`, () => {
      const req = makeReq({ httpMethod: 1 }, "http://fold.sh", {
        status: 500,
      });
      expect(req.fresh).toBe(false);
    });
    test(`requests with status 502 are not fresh`, () => {
      const req = makeReq({ httpMethod: 1 }, "http://fold.sh", {
        status: 502,
      });
      expect(req.fresh).toBe(false);
    });
    test(`if etags match it should be fresh`, () => {
      const req = makeReq(
        {
          httpMethod: 1,
          headersMap: headerMap({
            "If-None-Match": ['"foo"'],
          }),
        },
        "http://fold.sh",
        {
          status: 200,
          headersMap: headerMap({
            ETag: ['"foo"'],
          }),
        }
      );
      expect(req.fresh).toBe(true);
    });
    test(`if etags don't match it should be fresh`, () => {
      const req = makeReq(
        {
          httpMethod: 1,
          headersMap: headerMap({
            "If-None-Match": ['"foo"'],
          }),
        },
        "http://fold.sh",
        {
          status: 200,
          headersMap: headerMap({
            ETag: ['"bar"'],
          }),
        }
      );
      expect(req.fresh).toBe(false);
    });
    test(`if not modified it should be stale`, () => {
      const req = makeReq(
        {
          httpMethod: 1,
          headersMap: headerMap({
            "If-Modified-Since": ["Sat, 01 Jan 2000 00:00:00 GMT"],
          }),
        },
        "http://fold.sh",
        {
          status: 200,
          headersMap: headerMap({
            "Last-Modified": ["Sat, 01 Jan 2000 00:00:00 GMT"],
          }),
        }
      );
      expect(req.fresh).toBe(true);
    });
    test(`if modified it should be stale`, () => {
      const req = makeReq(
        {
          httpMethod: 1,
          headersMap: headerMap({
            "If-Modified-Since": ["Sat, 01 Jan 2000 00:00:00 GMT"],
          }),
        },
        "http://fold.sh",
        {
          status: 200,
          headersMap: headerMap({
            "Last-Modified": ["Sat, 01 Jan 2000 01:00:00 GMT"],
          }),
        }
      );
      expect(req.fresh).toBe(false);
    });
  });
  test("stale is the opposite of fresh", () => {
    const req = makeReq({ httpMethod: 1 }, "http://fold.sh", { status: 201 });
    expect(req.stale).toBe(true);
  });
});

// Note that we're testing against the Request interface.
// The goal here is to make sure that the request wrapper
// implements said interface properly.
export function makeReq(
  params: Partial<FoldHTTPRequest.AsObject>,
  url?: string,
  resParams?: Partial<FoldHTTPResponse.AsObject>
): Request {
  params = merge(newReq(), params);
  const foldReq = foldReqFromObject(params as FoldHTTPRequest.AsObject);
  if (url) {
    updateWithURL(foldReq, url);
  }
  const req = new FoldHTTPRequestWrapper(foldReq);
  resParams = merge(newRes(), resParams ? resParams : {});
  const foldRes = foldResFromObject(resParams as FoldHTTPResponse.AsObject);
  const res = new FoldHTTPResponseWrapper(req, foldRes);
  foldRes.getHeadersMap().forEach((values, header) => {
    res.set(header, values.getValuesList().join(HEADER_SEPARATOR));
  });
  req.res = res;
  return req;
}
