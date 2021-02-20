import { TextDecoder, TextEncoder } from "util";
import { GrpcClient } from "./__tests__/grpc_client";
import { Request as ProtoRequest } from "../dist/proto/ingress_pb";
import {
  Route,
  HttpMethod as ProtoHttpMethod,
} from "../dist/proto/manifest_pb";

import { Service, Request, Response } from "./service";
import { mockLogger } from "./logging";

describe("Fold Service", () => {
  const sockAddr: string = "/tmp/fold-node-sdk-test.sock";
  process.env.FOLD_SOCK_ADDR = sockAddr;
  const service: Service = new Service("test");
  service.logger = mockLogger();
  service.get("/get", (req: Request, res: Response) => {
    res.statusCode = 200;
    res.body = req.body;
  });
  service.put("/put", (req: Request, res: Response) => {
    res.statusCode = req.body.status;
    res.body = {
      method: req.httpMethod,
      path: req.path,
      body: req.body,
    };
  });
  service.start();
  const client: GrpcClient = new GrpcClient(sockAddr);
  afterAll(() => {
    client.stop();
    // Hacky but I dont' want to expose this in the Service API.
    (service as any).grpcBackend.server.forceShutdown();
  });
  describe("GetManifest", () => {
    it("should return the correct manifest", async () => {
      const manifest = await client.getManifest();
      expect(manifest.getName()).toEqual("test");
      const eRoute1 = new Route();
      eRoute1.setHttpMethod(ProtoHttpMethod.GET);
      eRoute1.setPathSpec("/get");
      const eRoute2 = new Route();
      eRoute2.setHttpMethod(ProtoHttpMethod.PUT);
      eRoute2.setPathSpec("/put");
      let expectation = [eRoute1, eRoute2];

      manifest!.getRoutesList().forEach((route, i) => {
        expect(route.getHttpMethod()).toEqual(expectation[i].getHttpMethod());
        expect(route.getPathSpec()).toEqual(expectation[i].getPathSpec());
      });
    });
  });

  describe("DoRequest", () => {
    describe("for a GET request", () => {
      it("should have an empty object for a body", async () => {
        const req = new ProtoRequest();
        const manifest = await client.getManifest();
        const handler = manifest.getRoutesList()[0].getHandler();
        req.setHttpMethod(ProtoHttpMethod.GET);
        req.setHandler(handler);
        req.setPath("/get");
        req.setBody(new TextEncoder().encode(JSON.stringify({})));

        const res = await client.doRequest(req);
        const resBody = JSON.parse(
          new TextDecoder("utf-8").decode(res.getBody_asU8())
        );
        expect(res.getStatus()).toEqual(200);
        expect(resBody).toEqual({});
      });
    });
    describe("for a given request body in a PUT", () => {
      const cases = [
        {
          httpMethod: "PUT",
          path: "/put",
          body: { status: 200, msg: "foo" },
        },
        {
          httpMethod: "PUT",
          path: "/put",
          body: { status: 1234, msg: "bar" },
        },
        {
          httpMethod: "PUT",
          path: "/put",
          body: { status: 202, msg: { a: "nested message" } },
        },
        {
          httpMethod: "PUT",
          path: "/put",
          body: { status: 214, msg: ["an", "array", "message"] },
        },
      ];
      it.each(cases)(
        "the server should return the correct response",
        async ({ httpMethod, path, body }) => {
          const manifest = await client.getManifest();
          const handler = manifest.getRoutesList()[1].getHandler();
          const req = new ProtoRequest();
          req.setHttpMethod(ProtoHttpMethod.PUT);
          req.setHandler(handler);
          req.setPath(path);
          req.setBody(new TextEncoder().encode(JSON.stringify(body)));

          const res = await client.doRequest(req);
          expect(res.getStatus()).toEqual(body.status);
          const resBody = JSON.parse(
            new TextDecoder("utf-8").decode(res.getBody_asU8())
          );
          expect(resBody).toEqual({
            method: httpMethod,
            path: path,
            body: body,
          });
        }
      );
    });
  });
});
