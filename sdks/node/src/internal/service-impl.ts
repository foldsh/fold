import { getLogger, Logger } from "../logging";
import { HTTPMethod, Request, Response } from "../http";
import { Version } from "../version";
import { GrpcServer } from "./grpc-server";
import {
  Middleware,
  MiddlewareFn,
  HandlerFn,
  NextFn,
  MiddlewareT,
  MiddlewareObj,
} from "../middleware";
import { Service } from "../service";
import { FoldHTTPRequestWrapper } from "./req-wrapper";
import { FoldHTTPResponseWrapper } from "./res-wrapper";
import { RouteTable } from "./route-table";
import { RouteTree } from "./route-tree";
import { ManifestSpec } from "./utils";

export class ServiceImpl implements Service {
  public readonly type: MiddlewareT = MiddlewareT.SERVICE;

  public name: string;
  public version: Version;
  public settings: ServiceSettings;
  public logger: Logger;
  public router!: RouteTable;

  public parent?: Service;
  public socket?: string;

  private grpcServer!: GrpcServer;
  private running: boolean = false;
  routes: RouteTree;

  constructor(version?: Version) {
    this.name = process.env.FOLD_SERVICE_NAME!;
    this.version = version ? version : { major: 0, minor: 0, patch: 1 };
    this.settings = new ServiceSettings();
    this.logger = getLogger(this.name);
    this.routes = new RouteTree();
  }

  public start(): void {
    if (this.running) {
      throw new Error("The service is already running");
      return;
    }
    const socket = this.socket ? this.socket : process.env.FOLD_SOCK_ADDR!;
    const router = RouteTable.fromRouteTree(this.logger, this.routes);
    const manifest: ManifestSpec = {
      name: this.name,
      version: this.version,
      routes: router.routeManifest(),
    };
    this.logger.debug(manifest);
    this.grpcServer = new GrpcServer(this.logger, router, manifest, socket);
    this.grpcServer.serve();
    this.running = true;
  }

  public get(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.GET, route, handler);
  }

  public head(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.HEAD, route, handler);
  }

  public post(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.POST, route, handler);
  }

  public put(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.PUT, route, handler);
  }

  public delete(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.DELETE, route, handler);
  }

  public connect(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.CONNECT, route, handler);
  }

  public options(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.OPTIONS, route, handler);
  }

  public trace(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.TRACE, route, handler);
  }

  public patch(route: string, handler: HandlerFn): void {
    this.routes.addHandler(HTTPMethod.PATCH, route, handler);
  }

  public all(route: string, handler: HandlerFn): void {
    for (const method in HTTPMethod) {
      this.routes.addHandler(method as HTTPMethod, route, handler);
    }
  }

  // This is essentially the method for adding middlewares. It's annoying
  // fiddly but it gives a nice and fleixble API as a result.
  public use(use: Middleware | string, middleware?: Middleware): void {
    if (middleware === undefined && typeof use === "string") {
      throw new TypeError(
        "If you specify a route you must also pass a middleware."
      );
    }
    if (arguments.length === 2 && typeof use !== "string") {
      throw new TypeError("You must specify the route as the first parameter");
    }

    // Now we need to sort out which parameter is which.
    let mw: Middleware;
    let route: string;
    if (middleware) {
      // Then the first argument is the route and the second is the middleware.
      route = use as string;
      mw = middleware;
    } else {
      // Otherwise, we use the default route.
      route = "/";
      // And use the first argument as the middleware.
      // what type it has yet.
      mw = use as Middleware;
    }

    // Now we need to resolve the middleware to the specific type
    let mwObj: MiddlewareObj;
    if (typeof mw === "function") {
      // It's either a HandlerFn or a MiddlewareFn
      if (mw.length === 2) {
        mwObj = { handle: mw as HandlerFn, type: MiddlewareT.HANDLER_FN };
      } else if (mw.length === 3) {
        mwObj = { handle: mw as MiddlewareFn, type: MiddlewareT.MIDDLEWARE_FN };
      } else {
        throw new TypeError(
          "Middleware interface requires either 2 or 3 parameters."
        );
      }
    } else {
      // It's a service, so it and the rotuer already satisfy the
      // MiddlewareObj interface.
      const svc: ServiceImpl = mw as ServiceImpl;
      svc.parent = this;
      mwObj = svc;
    }
    // We're done, so we add it to the router.
    this.routes.addMiddleware(route, mwObj);
  }

  public handle(req: Request, res: Response, next: NextFn): void {
    this.router.handle(
      req as FoldHTTPRequestWrapper,
      res as FoldHTTPResponseWrapper,
      next
    );
  }

  public param() {
    throw new Error("not yet implemented");
  }

  public route() {
    throw new Error("not yet implemented");
  }

  public shutdown(): void {
    this.grpcServer.shutdown();
    this.running = false;
  }
}

class ServiceSettings {
  private readonly settings: Map<string, any>;

  constructor() {
    this.settings = new Map<string, any>();
  }

  public get(setting: string): any | undefined {
    return this.settings.get(setting);
  }
  public set(setting: string, value: any): void {
    this.settings.set(setting, value);
  }

  public enable(setting: string): void {
    this.settings.set(setting, true);
  }

  public enabled(setting: string): boolean {
    return !!this.settings.get(setting);
  }

  public disable(setting: string): void {
    this.settings.set(setting, false);
  }
}
