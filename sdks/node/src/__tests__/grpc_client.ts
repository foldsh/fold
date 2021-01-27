import { unlinkSync } from "fs";
import { credentials } from "@grpc/grpc-js";
import { FoldIngressClient } from "../../dist/proto/ingress_grpc_pb";
import { Manifest } from "../../dist/proto/manifest_pb";
import { ManifestReq, Request, Response } from "../../dist/proto/ingress_pb";

export class GrpcClient {
  private client: FoldIngressClient;
  private sockAddr: string;

  constructor(sockAddr: string) {
    this.sockAddr = sockAddr;
    this.client = new FoldIngressClient(
      `unix://${sockAddr}`,
      credentials.createInsecure()
    );
  }

  public stop(): void {
    this.client.close();
    unlinkSync(this.sockAddr);
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

  public async doRequest(req: Request): Promise<Response> {
    return new Promise((resolve, reject): void => {
      this.client.doRequest(req, (err, res?: Response) => {
        if (err) {
          return reject(err);
        }
        resolve(res!);
      });
    });
  }
}
