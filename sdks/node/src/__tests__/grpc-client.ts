import { credentials } from "@grpc/grpc-js";
import { FoldIngressClient } from "../../dist/proto/ingress_grpc_pb";
import { Manifest } from "../../dist/proto/manifest_pb";
import { ManifestReq } from "../../dist/proto/ingress_pb";
import { FoldHTTPRequest, FoldHTTPResponse } from "../../dist/proto/http_pb";

export class GrpcClient {
  private client: FoldIngressClient;

  constructor(sockAddr: string) {
    this.client = new FoldIngressClient(
      `unix://${sockAddr}`,
      credentials.createInsecure()
    );
  }

  public stop(): void {
    this.client.close();
  }

  public async getManifest(): Promise<Manifest> {
    return new Promise((resolve, reject): void => {
      this.client.getManifest(new ManifestReq(), (err, manifest?: Manifest) => {
        if (err) {
          return reject(err);
        }
        resolve(manifest!);
      });
    });
  }

  public async doRequest(req: FoldHTTPRequest): Promise<FoldHTTPResponse> {
    return new Promise((resolve, reject): void => {
      this.client.doRequest(req, (err, res?: FoldHTTPResponse) => {
        if (err) {
          return reject(err);
        }
        resolve(res!);
      });
    });
  }
}
