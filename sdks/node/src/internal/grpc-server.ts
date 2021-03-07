import {
  Server,
  credentials,
  ServerUnaryCall,
  sendUnaryData,
} from "@grpc/grpc-js";

import { Manifest } from "../../dist/proto/manifest_pb";
import { FoldHTTPRequest, FoldHTTPResponse } from "../../dist/proto/http_pb";
import {
  FoldIngressService,
  IFoldIngressServer,
} from "../../dist/proto/ingress_grpc_pb";
import { ManifestReq } from "../../dist/proto/ingress_pb";

import { Logger } from "../logging";

import { FoldHTTPRequestWrapper } from "./req-wrapper";
import { FoldHTTPResponseWrapper } from "./res-wrapper";
import { makeManifest, ManifestSpec } from "./utils";
import { RouteTable } from "./route-table";

/**
 * THe gRPC server implementation that backs a service.
 */
export class GrpcServer {
  public logger: Logger;
  public manifest: Manifest;

  private router: RouteTable;
  private server!: Server;
  private socket: string;

  constructor(
    logger: Logger,
    router: RouteTable,
    manifest: ManifestSpec,
    socket: string
  ) {
    this.router = router;
    this.manifest = makeManifest(manifest);
    this.logger = logger;
    this.socket = socket;
  }

  /**
   * This works by dispatching a
   */
  handleRequest(
    req: FoldHTTPRequestWrapper,
    callback: (res: FoldHTTPResponseWrapper) => void
  ) {
    let res = new FoldHTTPResponseWrapper(req, new FoldHTTPResponse());
    req.res = res;
    this.router.handle(req, res, () => {
      callback(res);
    });
  }

  serve(): void {
    const socketAddr = `unix://${this.socket}`;
    this.logger.debug(`starting server on socket ${socketAddr}`);
    this.server = new Server();
    this.server.addService(FoldIngressService, newFoldIngressServer(this));
    this.server.bindAsync(
      socketAddr,
      credentials.createInsecure() as any,
      (err: Error | null, bindPort: number) => {
        if (err) {
          throw err;
        }
        this.logger.debug(
          `binding server to socket ${socketAddr} with port ${bindPort}`
        );
        this.server.start();
      }
    );
  }

  shutdown(): void {
    this.server.forceShutdown();
  }
}

function newFoldIngressServer(backend: GrpcServer): IFoldIngressServer {
  return {
    getManifest(
      _: ServerUnaryCall<ManifestReq, Manifest>,
      callback: sendUnaryData<Manifest>
    ): void {
      backend.logger.debug(`retrieving service manifest`);
      callback(null, backend.manifest);
    },

    doRequest(
      call: ServerUnaryCall<FoldHTTPRequest, FoldHTTPResponse>,
      callback: sendUnaryData<FoldHTTPResponse>
    ): void {
      let request = new FoldHTTPRequestWrapper(call.request);
      backend.logger.debug(
        `handling request for ${request.method} ${request.path}`
      );
      backend.handleRequest(request, (res: FoldHTTPResponseWrapper) => {
        callback(null, res.foldHTTPResponse);
      });
    },
  };
}
