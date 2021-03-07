import { URL } from "url";

import {
  FoldHTTPProto,
  FoldHTTPRequest,
  FoldHTTPResponse,
  StringArray,
} from "../../dist/proto/http_pb";

export function newRes(): FoldHTTPResponse.AsObject {
  return ({
    headersMap: [],
  } as unknown) as FoldHTTPResponse.AsObject;
}

export function foldResFromObject(
  obj: FoldHTTPResponse.AsObject
): FoldHTTPResponse {
  const res = new FoldHTTPResponse();
  res.setStatus(obj.status);
  res.setBody(obj.body);
  obj.headersMap.forEach(([key, values]) => {
    res.getHeadersMap().set(key, stringArrayFromObject(values));
  });
  return res;
}

export function paramsMap(map: {
  [key: string]: string;
}): Array<[string, string]> {
  let results: Array<[string, string]> = [];
  for (let [key, values] of Object.entries(map)) {
    results.push([key, values]);
  }
  return results;
}

export function queryMap(map: {
  [key: string]: string[];
}): Array<[string, StringArray.AsObject]> {
  let results: Array<[string, StringArray.AsObject]> = [];
  for (let [key, values] of Object.entries(map)) {
    results.push([key, { valuesList: values }]);
  }
  return results;
}

export function headerMap(map: {
  [key: string]: string[];
}): Array<[string, StringArray.AsObject]> {
  return queryMap(map);
}

export function newReq(): FoldHTTPRequest.AsObject {
  return ({
    url: {},
    httpProto: {},
    headersMap: [],
    pathParamsMap: [],
    queryParamsMap: [],
  } as unknown) as FoldHTTPRequest.AsObject;
}

export function foldReqFromObject(
  obj: FoldHTTPRequest.AsObject
): FoldHTTPRequest {
  const req = new FoldHTTPRequest();
  req.setHttpMethod(obj.httpMethod);
  req.setHttpProto(foldProtoFromObject(obj.httpProto));
  req.setHost(obj.host);
  req.setRemoteAddr(obj.remoteAddr);
  req.setRequestUri(obj.requestUri);
  req.setContentLength(obj.contentLength);
  if (obj.body !== undefined) {
    const body = Buffer.from(obj.body as string, "utf-8");
    req.setBody(body);
    req
      .getHeadersMap()
      .set(
        "transfer-encoding",
        stringArrayFromObject({ valuesList: ["utf-8"] })
      );
    req
      .getHeadersMap()
      .set(
        "content-length",
        stringArrayFromObject({ valuesList: [String(body.length)] })
      );
  }
  obj.headersMap.forEach(([key, values]) => {
    req.getHeadersMap().set(key, stringArrayFromObject(values));
  });
  obj.pathParamsMap.forEach(([key, value]) => {
    req.getPathParamsMap().set(key, value);
  });
  obj.queryParamsMap.forEach(([key, values]) => {
    req.getQueryParamsMap().set(key, stringArrayFromObject(values));
  });
  req.setRoute(obj.route);
  return req;
}

export function updateWithURL(req: FoldHTTPRequest, url: string): void {
  const parsed = new URL(url);
  req.setPath(parsed.pathname.replace("?", ""));
  req.setRawQuery(parsed.search.replace("?", ""));
  req.setFragment(parsed.hash.replace("#", ""));
}

function foldProtoFromObject(obj?: FoldHTTPProto.AsObject): FoldHTTPProto {
  const proto = new FoldHTTPProto();
  if (!obj) {
    return proto;
  }
  proto.setProto(obj.proto);
  proto.setMajor(obj.major);
  proto.setMinor(obj.minor);
  return proto;
}

function stringArrayFromObject(obj: StringArray.AsObject): StringArray {
  const s = new StringArray();
  s.setValuesList(obj.valuesList);
  return s;
}
