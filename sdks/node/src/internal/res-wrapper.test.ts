import merge from "utils-merge";
import { FoldHTTPRequest, FoldHTTPResponse } from "../../dist/proto/http_pb";
import {
  newReq,
  foldReqFromObject,
  newRes,
  foldResFromObject,
  headerMap,
} from "../__tests__/proto-utils";
import { Response } from "../http";
import { FoldHTTPRequestWrapper } from "./req-wrapper";
import { FoldHTTPResponseWrapper } from "./res-wrapper";
import { HEADER_SEPARATOR } from "./utils";

describe("FoldHTTPResponseWrapper", () => {
  describe("res.append", () => {
    const res = makeRes({
      headersMap: headerMap({
        "Set-Cookie": ["id=a3fWa"],
        Link: ['<https://one.example.com>; rel="preconnect"'],
      }),
    });
    test("should be able to append to set-cookie", () => {
      res.append("Set-Cookie", "Expires=Wed, 21 Oct 2015 07:28:00 GMT");
      expect(res.get("Set-Cookie")).toEqual([
        "id=a3fWa",
        "Expires=Wed, 21 Oct 2015 07:28:00 GMT",
      ]);
    });
    test("should be able to append to set-cookie with array", () => {
      res.append("Set-Cookie", ["Secure"]);
      expect(res.get("Set-Cookie")).toEqual([
        "id=a3fWa",
        "Expires=Wed, 21 Oct 2015 07:28:00 GMT",
        "Secure",
      ]);
    });
    test("should be concatenate arrays that are not for set-cookie", () => {
      res.append("Cache-Control", ["no-cache", "no-store"]);
      expect(res.get("Cache-Control")).toEqual("no-cache, no-store");
    });
    test("should be able to append to other headers", () => {
      res.append("Link", '<https://two.example.com>; rel="preconnect"');
      expect(res.get("Link")).toEqual(
        '<https://one.example.com>; rel="preconnect", <https://two.example.com>; rel="preconnect"'
      );
      res.append("Link", '<https://three.example.com>; rel="preconnect"');
      expect(res.get("Link")).toEqual(
        '<https://one.example.com>; rel="preconnect", <https://two.example.com>; rel="preconnect", <https://three.example.com>; rel="preconnect"'
      );
    });
  });

  describe("res.clearCookie", () => {
    test("should set the cookie value to empty after settign it", () => {
      const res = makeRes({});
      const lax = "lax" as const;
      const options = { httpOnly: true, sameSite: lax };
      res.clearCookie("fold", options);
      expect(res.get("Set-Cookie")).toEqual(
        "fold=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; SameSite=Lax"
      );
    });
  });

  describe("res.cookie", () => {
    test("should set a cookie value with default Path", () => {
      const res = makeRes({});
      const options = { httpOnly: true };
      res.cookie("fold", "value", options);
      expect(res.get("Set-Cookie")).toEqual("fold=value; Path=/; HttpOnly");
    });
    test("should sign a cookie if required", () => {
      const res = makeRes({});
      const req = (res as FoldHTTPResponseWrapper).req;
      req.secret = "secret";
      res.cookie("fold", "value", { signed: true });
      expect(res.get("Set-Cookie")).toEqual(
        "fold=s%3Avalue.UOA%2BvmW%2BmLuL8RuiyJLVTAeayisNOwFidpxtdXolQ08; Path=/"
      );
    });
    test("should thrown an error if asked to sign with no secret", () => {
      const res = makeRes({});
      expect(() => {
        res.cookie("fold", "value", { signed: true });
      }).toThrowError(Error);
    });
    test("should stringify json value, prefix with j:, and url encode it", () => {
      const res = makeRes({});
      res.cookie("fold", { hello: "world" }, {});
      expect(res.get("Set-Cookie")).toEqual(
        `fold=j%3A%7B%22hello%22%3A%22world%22%7D; Path=/`
      );
    });
    test("should encode cookie if required", () => {
      const res = makeRes({});
      const b64e = (value: string): string => {
        return Buffer.from(value).toString("base64");
      };
      const encodedValue = b64e("value");
      res.cookie("fold", "value", { encode: b64e });
      expect(res.get("Set-Cookie")).toEqual(`fold=${encodedValue}; Path=/`);
    });
    test("should set max age correctly", () => {
      jest
        .spyOn(global.Date, "now")
        .mockImplementation(() => new Date(1).valueOf());
      const res = makeRes({});
      res.cookie("fold", "value", { maxAge: 1000000 });
      const expires = new Date(1 + 1000000);
      expect(res.get("Set-Cookie")).toEqual(
        `fold=value; Max-Age=1000; Path=/; Expires=${expires.toUTCString()}`
      );
    });
    test("should allow multiple calls", () => {
      const res = makeRes({});
      res.cookie("fold", "value");
      res.cookie("foo", "?");
      res.cookie("bar", "@");
      expect(res.get("Set-Cookie")).toEqual(
        "fold=value; Path=/; foo=%3F; Path=/; bar=%40; Path=/"
      );
    });
  });

  describe("res.set and res.get", () => {
    const res = makeRes({});
    test("should be able to retrieve it", () => {
      res.set("Content-Length", "1024");
      expect(res.get("Content-Length")).toEqual("1024");
    });
    test("should only be able to set 'set-cookie' to an array", () => {
      expect(() => {
        res.set("Content-Length", ["foo", "bar"]);
      }).toThrowError(TypeError);
      const cookies = ["id=a3fWa", "Expires=Wed, 21 Oct 2015 07:28:00 GMT"];
      res.set("Set-Cookie", cookies);
      expect(res.get("Set-Cookie")).toEqual(cookies);
    });
    test("charset should be added to content-type if not present", () => {
      res.set("content-type", "application/json");
      expect(res.get("content-type")).toEqual(
        "application/json; charset=utf-8"
      );
    });
    test("charset should not be added to content-type if present", () => {
      res.set("content-type", "application/json; charset=utf-8");
      expect(res.get("content-type")).toEqual(
        "application/json; charset=utf-8"
      );
    });
    test("should be able to add multiple headers at once", () => {
      res.set({ "Content-Type": "application/json", "Content-Length": "1024" });
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
      expect(res.get("Content-Length")).toEqual("1024");
    });
  });

  describe("res.links", () => {
    const res = makeRes({});
    test("should create links with none set", () => {
      res.links({ foo: "www.foo.com", bar: "www.bar.com" });
      expect(res.get("Link")).toEqual(
        '<www.foo.com>; rel="foo", <www.bar.com>; rel="bar"'
      );
    });
    test("should add to existing links", () => {
      res.links({ baz: "www.baz.com", fold: "www.fold.sh" });
      expect(res.get("Link")).toEqual(
        '<www.foo.com>; rel="foo", <www.bar.com>; rel="bar", <www.baz.com>; rel="baz", <www.fold.sh>; rel="fold"'
      );
    });
  });

  describe("res.location", () => {
    const res = makeRes(
      {},
      { headersMap: headerMap({ Referrer: ["www.fold.sh"] }) }
    );
    test("should return the referrer if set", () => {
      res.location("back");
      expect(res.get("Location")).toEqual("www.fold.sh");
    });
    test("should return default referrer of '/'", () => {
      const res = makeRes({});
      res.location("back");
      expect(res.get("Location")).toEqual("/");
    });
    test("should location as specified if not 'back'", () => {
      res.location("notback.fold.sh");
      expect(res.get("Location")).toEqual("notback.fold.sh");
    });
  });

  describe("res.redirect", () => {
    test("should redirect with default 302 status", () => {
      const res = makeRes({});
      res.redirect("www.fold.sh");
      expect(res.get("Location")).toEqual("www.fold.sh");

      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(foldRes.getStatus()).toEqual(302);
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual(`{"title":"Redirecting to www.fold.sh"}`);
    });
    test("should redirect with specified status", () => {
      const res = makeRes({});
      res.redirect("www.fold.sh", 304);
      expect(res.get("Location")).toEqual("www.fold.sh");
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(foldRes.getStatus()).toEqual(304);
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual(`{"title":"Redirecting to www.fold.sh"}`);
    });
  });

  describe("res.send", () => {
    test("should send a text body when given a string", () => {
      const res = makeRes({});
      res.send("request body");
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(res.get("Content-Type")).toEqual("text/plain; charset=utf-8");
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual("request body");
    });
    test("should send a json body when given a number", () => {
      const res = makeRes({});
      res.send(1);
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual("1");
    });
    test("should send a json body when given a boolean", () => {
      const res = makeRes({});
      res.send(true);
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual("true");
    });
    test("should send a json body when given a object", () => {
      const res = makeRes({});
      res.send({ foo: "bar" });
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual(`{"foo":"bar"}`);
    });
    test("should send a binary body when given a buffer", () => {
      const res = makeRes({});
      res.send(Buffer.from([123, 125]));
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(res.get("Content-Type")).toEqual("application/octet-stream");
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual("{}");
    });
  });

  describe("res.json", () => {
    test("should send a json body", () => {
      const res = makeRes({});
      res.json({ foo: "bar" });
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
      expect(body).toEqual(`{"foo":"bar"}`);
    });
  });

  describe("res.sendStatus", () => {
    test("should set the status and send the appropriate message in the body", () => {
      const res = makeRes({});
      res.sendStatus(200);
      const foldRes = (res as FoldHTTPResponseWrapper).foldHTTPResponse;
      expect(foldRes.getStatus()).toEqual(200);
      expect(res.get("Content-Type")).toEqual("text/plain; charset=utf-8");
      const body = Buffer.from(foldRes.getBody_asU8()).toString("utf-8");
      expect(body).toEqual(`OK`);
    });
  });

  describe("res.status", () => {
    test("should set the response status code", () => {
      const res = makeRes({});
      res.status(204);
      expect(res.statusCode).toEqual(204);
      expect(
        (res as FoldHTTPResponseWrapper).foldHTTPResponse.getStatus()
      ).toEqual(204);
    });
  });

  describe("res.type", () => {
    test("should set the json content type appropriately", () => {
      const res = makeRes({});
      res.type("json");
      expect(res.get("Content-Type")).toEqual(
        "application/json; charset=utf-8"
      );
    });
    test("should set the text/plain content type appropriately", () => {
      const res = makeRes({});
      res.type("txt");
      expect(res.get("Content-Type")).toEqual("text/plain; charset=utf-8");
    });
  });

  describe("res.vary", () => {
    const res = makeRes({});
    test("should set the vary header when none is set", () => {
      res.vary("User-Agent");
      expect(res.get("Vary")).toEqual("User-Agent");
    });
    test("should add to the vary header when it is already set", () => {
      res.vary("Referrer");
      expect(res.get("Vary")).toEqual("User-Agent, Referrer");
    });
  });
});

export function makeRes(
  params: Partial<FoldHTTPResponse.AsObject>,
  reqParams?: Partial<FoldHTTPRequest.AsObject>
): Response {
  // Create the request
  reqParams = merge(newReq(), reqParams ? reqParams : {});
  const foldReq = foldReqFromObject(reqParams as FoldHTTPRequest.AsObject);
  const req = new FoldHTTPRequestWrapper(foldReq);

  // Create the response
  params = merge(newRes(), params);
  const foldRes = foldResFromObject(params as FoldHTTPResponse.AsObject);
  const res = new FoldHTTPResponseWrapper(req, foldRes);
  foldRes.getHeadersMap().forEach((values, header) => {
    const vals =
      header.toLowerCase() === "set-cookie"
        ? values.getValuesList()
        : values.getValuesList().join(HEADER_SEPARATOR);
    res.set(header, vals);
  });
  return res;
}
