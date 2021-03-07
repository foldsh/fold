import { RouteTable } from "./route-table";
import { RouteTree } from "./route-tree";
import { mockLogger } from "../logging";
import { HandlerFn, MiddlewareObj, MiddlewareT } from "../middleware";
import { makeRes } from "./res-wrapper.test";

describe("RouteTable", () => {
  test("it should call a basic handler", () => {
    const handler = jest.fn((_req, res) => res.send("done"));
    const done = jest.fn();
    const table = buildRouteTable({ "/": { GET: handler } });
    table.handle({ method: "GET", route: "/" } as any, makeRes({}), done);
    expect(handler).toHaveBeenCalledTimes(1);
    expect(done).toHaveBeenCalledTimes(1);
  });
  test("it should call a registered middleware before the handler", () => {
    const calls: number[] = [];
    const middleware = jest.fn((_req, _res, next) => {
      calls.push(1);
      next();
    });
    const handler = jest.fn((_req, res) => {
      calls.push(2);
      res.send("done");
    });
    const done = jest.fn(() => {
      calls.push(3);
    });
    const table = buildRouteTable(
      { "/": { GET: handler } },
      { "/": [mware(middleware)] }
    );
    table.handle({ method: "GET", route: "/" } as any, makeRes({}), done);
    expect(middleware).toHaveBeenCalledTimes(1);
    expect(handler).toHaveBeenCalledTimes(1);
    expect(done).toHaveBeenCalledTimes(1);
    expect(calls).toEqual([1, 2, 3]);
  });
  test("it should call all middlwares matching the start of a path", () => {
    const calls: number[] = [];
    const m1 = jest.fn((_req, _res, next) => {
      calls.push(1);
      next();
    });
    const m2 = jest.fn((_req, _res, next) => {
      calls.push(2);
      next();
    });
    const m3 = jest.fn((_req, _res, next) => {
      calls.push(3);
      next();
    });
    const m4 = jest.fn((_req, _res, next) => {
      calls.push(4);
      next();
    });
    const m5 = jest.fn((_req, _res, next) => {
      calls.push(5);
      next();
    });
    const m6 = jest.fn((_req, _res, next) => {
      calls.push(6);
      next();
    });
    const handler = jest.fn((_req, res) => {
      calls.push(7);
      res.send("done");
    });
    const done = jest.fn(() => {
      calls.push(8);
    });
    const table = buildRouteTable(
      { "/foo/bar": { GET: handler } },
      {
        "/": [mware(m1), mware(m2)],
        "/foo": [mware(m3)],
        "/foo/bar": [mware(m4)],
        "/foo/bar/baz": [mware(m6)],
        "/blah": [mware(m5)],
      }
    );
    table.handle(
      { method: "GET", route: "/foo/bar" } as any,
      makeRes({}),
      done
    );
    expect(m1).toHaveBeenCalledTimes(1);
    expect(m2).toHaveBeenCalledTimes(1);
    expect(m3).toHaveBeenCalledTimes(1);
    expect(m4).toHaveBeenCalledTimes(1);
    expect(m5).toHaveBeenCalledTimes(0);
    expect(m6).toHaveBeenCalledTimes(0);
    expect(handler).toHaveBeenCalledTimes(1);
    expect(done).toHaveBeenCalledTimes(1);
    expect(calls).toEqual([1, 2, 3, 4, 7, 8]);
  });

  test("it should call middlwares even when no handler route matches", () => {
    const calls: number[] = [];
    const m1 = jest.fn((_req, _res, next) => {
      calls.push(1);
      next();
    });
    const m2 = jest.fn((_req, _res, next) => {
      calls.push(2);
      next();
    });
    const done = jest.fn(() => {
      calls.push(3);
    });
    const table = buildRouteTable(
      {},
      {
        "/": [mware(m1), mware(m2)],
      }
    );
    table.handle(
      { method: "GET", route: "/metrics" } as any,
      makeRes({}),
      done
    );
    expect(m1).toHaveBeenCalledTimes(1);
    expect(m2).toHaveBeenCalledTimes(1);
    expect(done).toHaveBeenCalledTimes(1);
    expect(calls).toEqual([1, 2, 3]);
  });
  test("it should call done if there are no matching routes", () => {
    const calls: number[] = [];
    const done = jest.fn(() => {
      calls.push(1);
    });
    const table = buildRouteTable({}, {});
    table.handle({ method: "GET", route: "/" } as any, makeRes({}), done);
    expect(done).toHaveBeenCalledTimes(1);
    expect(calls).toEqual([1]);
  });
});

function buildRouteTable(
  handlers: { [key: string]: { [key: string]: Function } },
  middlewares?: { [key: string]: MiddlewareObj[] }
): RouteTable {
  const tree = new RouteTree();
  for (const [route, methods] of Object.entries(handlers)) {
    for (const [method, handler] of Object.entries(methods)) {
      tree.addHandler(method, route, handler as HandlerFn);
    }
  }
  if (middlewares) {
    for (const [route, middlewareObjs] of Object.entries(middlewares)) {
      for (const middleware of middlewareObjs) {
        tree.addMiddleware(route, middleware);
      }
    }
  }
  return RouteTable.fromRouteTree(mockLogger(), tree);
}

function mware(f: Function): MiddlewareObj {
  return {
    type: MiddlewareT.MIDDLEWARE_FN,
    handle: f as any,
  };
}
