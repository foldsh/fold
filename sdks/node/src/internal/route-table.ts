import {
  MiddlewareObj,
  MiddlewareT,
  HandlerFn,
  MiddlewareFn,
  NextFn,
} from "../middleware";
import { Request, Response } from "../http";

import { FoldHTTPRequestWrapper } from "./req-wrapper";
import { FoldHTTPResponseWrapper } from "./res-wrapper";
import { RouteSpec } from "./utils";
import { Logger } from "src/logging";
import { RouteTree, HandlerTable, MiddlewareTable } from "./route-tree";

/**
 * This isn't really a router in the traditional sense of the word.
 * All the path parsing etc happens in the runtime and in the SDK we
 * just get given the route that was matched, all the path params,
 * etc.
 *
 * This means we can take a bit of a shortcut here and just build
 * a very simple flat structure which we can then use to look up
 * the applicable middlewares and handlers from.
 *
 * We still use a more traditional tree structure to allow the developer
 * to define their application how they wish, bu we flatten it out on
 * start up.
 */
export class RouteTable implements MiddlewareObj {
  public readonly type: MiddlewareT = MiddlewareT.ROUTER;
  private logger: Logger;

  private handlerTable: HandlerTable;
  private middlewareTable: MiddlewareTable;

  constructor(
    logger: Logger,
    handlers: HandlerTable,
    middlewares: MiddlewareTable
  ) {
    this.logger = logger;
    this.handlerTable = handlers;
    this.middlewareTable = middlewares;
  }

  public handle(req: Request, res: Response, done: NextFn): void {
    const logger = this.logger;
    const foldReq = req as FoldHTTPRequestWrapper;
    const foldRes = res as FoldHTTPResponseWrapper;
    const route = foldReq.route;

    // First we get the full list of middlewares that will be applicable
    // to this route. This defaults to an empty list of there are none.
    const stack = this.applicableMiddlewares(route);

    // Now we need to figure out if there is an actual handler to apply.
    let handler: HandlerFn | undefined;
    const methodsOnRoute = this.handlerTable.get(route);
    if (methodsOnRoute) {
      handler = methodsOnRoute.get(req.method!);
    }
    // If there is, we push it on the stack of middleware like objects.
    if (handler) {
      stack.push({ handle: handler, type: MiddlewareT.HANDLER_FN });
    }
    // Before proceeding, we need to set up the call to done for when
    // the response is finished.
    res.on("finish", () => {
      logger.debug("finished handling request");
      done();
    });
    // If the stack is emtpy then it's a 404.
    if (stack.length === 0) {
      foldRes.status(404);
      foldRes.json({ title: "Page not found", detail: route });
      return;
    }
    // Ok then, we're ready to kick off the chain of handlers.
    next();

    // Our imlementation here can be relatively simple, due to the nature of
    // our middleware table. Because we've already resolved every route to
    // a flat list of handlers, we know that there are not any services or
    // routers in the list.
    //
    // We therefore just need to work our way through the list, passing in
    // next when reuired and calling done ourselves if we find there is
    // nothing left to do.
    function next(err?: any): void {
      const mw = stack.shift();
      if (err) {
        foldRes.status(500);
        foldRes.json({ title: "Internal server error", detail: err });
        return;
      }
      if (!mw) {
        // Nothing left that wants to handle this request.
        done();
        return;
      }
      switch (mw.type) {
        case MiddlewareT.HANDLER_FN:
          (mw.handle as HandlerFn)(foldReq, foldRes);
          return;
        case MiddlewareT.MIDDLEWARE_FN:
          (mw.handle as MiddlewareFn)(foldReq, foldRes, next);
          return;
        default:
          throw new Error(
            "This is a bug, a fully resolved route table should not contain any Services or Routers."
          );
      }
    }
  }

  private applicableMiddlewares(route: string): MiddlewareObj[] {
    // Get the route part of the handlerId
    let applicable: MiddlewareObj[] = [];
    for (let [keyRoute, middlewares] of this.middlewareTable.entries()) {
      if (route.startsWith(keyRoute)) {
        applicable = applicable.concat(middlewares);
      }
    }
    return applicable;
  }

  public routeManifest(): RouteSpec[] {
    const routes: RouteSpec[] = [];
    for (const [route, methods] of this.handlerTable.entries()) {
      for (const method of methods.keys())
        routes.push({
          method: method,
          handler: `${method} ${route}`,
          route: route,
        });
    }
    return routes;
  }

  public static fromRouteTree(
    logger: Logger,
    routeTree: RouteTree
  ): RouteTable {
    const flattenedTree = routeTree.flatten();
    const router = new RouteTable(
      logger,
      flattenedTree.handlers,
      flattenedTree.middlewares
    );
    // First up
    return router;
  }
}
