import { HandlerFn, MiddlewareObj, MiddlewareT } from "../middleware";
import { ServiceImpl } from "./service-impl";

export type HandlerTable = Map<string, Map<string, HandlerFn>>;
export type MiddlewareTable = Map<string, MiddlewareObj[]>;

const TRAILING_SLASH = /\/$/;

/**
 * The RouteTree exists purely to build up a tree like structure
 * of routes that is defined by the developer.
 * It is then used as the input to construct a full route table.
 */
export class RouteTree {
  handlers: HandlerTable;
  middlewares: MiddlewareTable;
  constructor() {
    this.handlers = new Map();
    this.middlewares = new Map();
  }

  public addHandler(method: string, route: string, handler: HandlerFn): void {
    const r = this.validateAndCleanRoute(route);
    let handlers = this.handlers.get(r);
    if (!handlers) {
      this.handlers.set(r, new Map());
      handlers = this.handlers.get(r);
    }
    handlers!.set(method, handler);
  }

  public addMiddleware(route: string, middleware: MiddlewareObj) {
    const r = this.validateAndCleanRoute(route);
    let middlewares = this.middlewares.get(r);
    if (!middlewares) {
      this.middlewares.set(r, []);
      middlewares = this.middlewares.get(r);
    }
    middlewares!.push(middleware);
  }

  private validateAndCleanRoute(path: string): string {
    if (!path || path[0] !== "/") {
      throw new Error("Path must be set and must start with a /");
    }
    // Remove any double slashes we've ended up with through nesting.
    let cleaned = path.replace("//", "/");
    // Remove a trailing slash, unless it's the root route '/'
    if (cleaned !== "/") {
      cleaned = cleaned.replace(TRAILING_SLASH, "");
    }
    return cleaned;
  }

  // Building the whole thing again is ineffecient, but it just
  // happens once on start up and it makes the algorithm much simpler
  // I think.
  flatten(): RouteTree {
    const tree = new RouteTree();
    // In order to flatten we first take the handlers from the current
    // level and add them with the routePrefix.
    copyRoutes(this, tree);
    // Then we need to through the middlewares, appending them to
    // the list of middlewares at this level, including the routePrefix.
    // If we come across a service, then we need to get its route tree
    // and flatten that too. Then we go through it and add all its
    // handlers and middlewares back to the current tree.
    for (const [route, middlewares] of this.middlewares.entries()) {
      for (const middleware of middlewares) {
        // const prefixedRoute = `${routePrefix}${route}`;
        switch (middleware.type) {
          case MiddlewareT.MIDDLEWARE_FN:
            // This case is easy, we just add it to the new tree.
            tree.addMiddleware(route, middleware);
            break;
          case MiddlewareT.SERVICE:
            // This is the recursive case, it's a bit tricker.
            // First we need to flatten it.
            const flattenedService = (middleware as ServiceImpl).routes.flatten();
            // const flattenedService = flatten(
            //   (middleware as ServiceImpl).routes
            // );
            // Then we need to go through its routes and copy them in.
            copyRoutes(flattenedService, tree, route);
            // And now we need to copy in the middlewares. As we have
            // just flattened the tree, there is no need to check for
            // services. We know that the middlewares on the flattened
            // routes will only contain plain middlewares.
            copyMiddlewares(flattenedService, tree, route);
            break;
          default:
            // If it's not a service or middleware fn then something has gone
            // horribly wrong so a crash is appropriate.
            throw new Error(
              "This is a bug, a fully resolved route table should not contain any Services or Routers."
            );
        }
      }
    }
    return tree;
  }
}

function copyRoutes(
  from: RouteTree,
  to: RouteTree,
  routePrefix: string = ""
): void {
  for (const [route, methods] of from.handlers.entries()) {
    for (const [method, handler] of methods.entries()) {
      to.addHandler(method, `${routePrefix}${route}`, handler);
    }
  }
}

function copyMiddlewares(
  from: RouteTree,
  to: RouteTree,
  routePrefix: string
): void {
  for (const [route, middlewares] of from.middlewares.entries()) {
    for (const middleware of middlewares) {
      to.addMiddleware(`${routePrefix}${route}`, middleware);
    }
  }
}
