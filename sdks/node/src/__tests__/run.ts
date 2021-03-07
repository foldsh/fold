import { v4 as uuidv4 } from "uuid";

import { Service } from "../service";
import { ServiceImpl } from "../internal/service-impl";
import {
  FoldHTTPMethod,
  FoldHTTPMethodMap,
  FoldHTTPResponse,
} from "../../dist/proto/http_pb";
import { GrpcClient } from "./grpc-client";
import { Manifest } from "../../dist/proto/manifest_pb";
import {
  foldReqFromObject,
  updateWithURL,
  headerMap,
  newReq,
  paramsMap,
  queryMap,
} from "./proto-utils";
import { joinHeader, makeManifest, ManifestSpec } from "../internal/utils";

export default function run(svc: Service): RequestRunner {
  return new RequestRunner(svc as ServiceImpl);
}

export interface TestRequest {
  route: string;
  params?: { [key: string]: string };
  query?: { [key: string]: string[] };
  headers?: { [key: string]: string[] };
  body?: any;
  url?: string;
}

export class RequestRunner {
  private client!: GrpcClient;
  private svc: ServiceImpl;
  private _manifest!: Manifest;
  private response!: FoldHTTPResponse;
  private actions: (() => Promise<void>)[];

  constructor(svc: ServiceImpl) {
    this.svc = svc;
    this.actions = [];
  }

  public get(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.GET, request);
  }

  public head(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.HEAD, request);
  }

  public post(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.POST, request);
  }

  public put(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.PUT, request);
  }

  public delete(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.DELETE, request);
  }

  public connect(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.CONNECT, request);
  }

  public options(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.OPTIONS, request);
  }

  public trace(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.TRACE, request);
  }

  public patch(request: TestRequest): RequestRunner {
    return this.handleRequest(FoldHTTPMethod.PATCH, request);
  }

  public manifest(): RequestRunner {
    this.actions.push(async () => {
      const manifest = await this.client.getManifest();
      this._manifest = manifest;
    });
    return this;
  }

  public expectStatus(status: number): RequestRunner {
    this.actions.push(async () => {
      expect(this.response.getStatus()).toEqual(status);
    });
    return this;
  }

  public expectBody(body: any): RequestRunner {
    this.actions.push(async () => {
      const resBody = Buffer.from(this.response.getBody_asU8()).toString(
        "utf-8"
      );
      switch (typeof body) {
        case "object":
          expect(JSON.parse(resBody)).toEqual(body);
          break;
        default:
          expect(resBody).toEqual(body);
          break;
      }
    });
    return this;
  }

  public expectHeader(header: string, value: any): RequestRunner {
    header = header.toLowerCase();
    this.actions.push(async () => {
      const headers = this.response.getHeadersMap();
      const resHeader = headers.get(header);
      if (resHeader) {
        expect(joinHeader(header, resHeader.getValuesList())).toEqual(value);
      } else {
        if (value === undefined) {
          expect(true).toBe(true);
        } else {
          fail(`header ${header} was not set on the response`);
        }
      }
    });
    return this;
  }

  public expectManifest(manifest: ManifestSpec): RequestRunner {
    this.actions.push(async () => {
      const actual = this._manifest;
      const expectation = makeManifest(manifest);
      expect(actual.toObject()).toEqual(expectation.toObject());
    });
    return this;
  }

  public async done(): Promise<void> {
    const socket = `/tmp/fold.${uuidv4()}.sock`;
    this.svc.socket = socket;
    try {
      this.svc.start();
      this.client = new GrpcClient(socket);
      for (const action of this.actions) {
        await action();
      }
    } finally {
      this.client.stop();
      this.svc.shutdown();
    }
  }

  private handleRequest(method: number, request: TestRequest): RequestRunner {
    const params = newReq();
    const foldMethod = method as FoldHTTPMethodMap[keyof FoldHTTPMethodMap];
    params.httpMethod = foldMethod;
    params.route = request.route;
    if (request.headers) {
      params.headersMap = headerMap(request.headers);
    }
    if (request.params) {
      params.pathParamsMap = paramsMap(request.params);
    }
    if (request.query) {
      params.queryParamsMap = queryMap(request.query);
    }
    if (request.body) {
      const body = request.body;
      switch (typeof body) {
        case "object":
          params.body = Buffer.from(JSON.stringify(body), "utf-8");
          break;
        default:
          params.body = Buffer.from(body, "utf-8");
      }
    }
    const req = foldReqFromObject(params);
    if (request.url) {
      updateWithURL(req, request.url);
    }
    this.actions.push(async () => {
      const res = await this.client.doRequest(req);
      this.response = res;
    });
    return this;
  }
}
