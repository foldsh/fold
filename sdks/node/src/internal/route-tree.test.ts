import { MiddlewareObj, MiddlewareT } from "../middleware";
import { RouteTree } from "./route-tree";
import { ServiceImpl } from "./service-impl";

describe("RouteTree", () => {
  test("it should flatten a simple tree", () => {
    const tree = new RouteTree();
    tree.addHandler("GET", "/foo", jest.fn());
    tree.addHandler("POST", "/foo", jest.fn());
    tree.addHandler("PUT", "/foo", jest.fn());

    tree.addHandler("GET", "/bar", jest.fn());
    tree.addHandler("POST", "/bar", jest.fn());
    tree.addHandler("PUT", "/bar", jest.fn());

    tree.addHandler("GET", "/baz/foo", jest.fn());
    tree.addHandler("POST", "/baz/foo", jest.fn());
    tree.addHandler("PUT", "/baz/foo", jest.fn());

    tree.addMiddleware("/", mware(jest.fn()));
    tree.addMiddleware("/", mware(jest.fn()));
    tree.addMiddleware("/", mware(jest.fn()));

    tree.addMiddleware("/foo", mware(jest.fn()));
    tree.addMiddleware("/foo", mware(jest.fn()));
    tree.addMiddleware("/foo", mware(jest.fn()));

    const flattened = tree.flatten();

    expect(flattened).toEqual(tree);
  });

  test("it should flatten a nested tree", () => {
    // Top level tree.
    const topLevel = new RouteTree();
    topLevel.addHandler("GET", "/foo", jest.fn());
    topLevel.addHandler("GET", "/bar", jest.fn());
    topLevel.addHandler("GET", "/baz/foo", jest.fn());

    topLevel.addMiddleware("/", mware(jest.fn()));

    // First branch of top level.
    const firstBranch: ServiceImpl = {
      type: MiddlewareT.SERVICE,
      routes: new RouteTree(),
    } as ServiceImpl;
    topLevel.addMiddleware("/baz", firstBranch);
    firstBranch.routes.addHandler("POST", "/foo", jest.fn());
    firstBranch.routes.addHandler("PUT", "/foo", jest.fn());

    // Second branch of top level
    const secondBranch: ServiceImpl = {
      type: MiddlewareT.SERVICE,
      routes: new RouteTree(),
    } as ServiceImpl;
    topLevel.addMiddleware("/foo", secondBranch);

    secondBranch.routes.addHandler("POST", "/:bar", jest.fn());
    secondBranch.routes.addHandler("PUT", "/:baz", jest.fn());

    secondBranch.routes.addMiddleware("/", mware(jest.fn()));
    secondBranch.routes.addMiddleware("/foo", mware(jest.fn()));

    // First branch of second branch
    const firstBranchOfSecondBranch: ServiceImpl = {
      type: MiddlewareT.SERVICE,
      routes: new RouteTree(),
    } as ServiceImpl;
    secondBranch.routes.addMiddleware("/:blah", firstBranchOfSecondBranch);
    firstBranchOfSecondBranch.routes.addHandler("GET", "/foo/blah", jest.fn());
    firstBranchOfSecondBranch.routes.addHandler("PUT", "/baz", jest.fn());
    firstBranchOfSecondBranch.routes.addMiddleware("/", mware(jest.fn()));
    firstBranchOfSecondBranch.routes.addMiddleware(
      "/metrics",
      mware(jest.fn())
    );

    // firstBranchOfSecondBranch.routes.addMiddleware("/foo", mware(jest.fn()));

    // Expectation
    const expectation = new RouteTree();
    // Top level
    expectation.addHandler("GET", "/foo", jest.fn());
    expectation.addHandler("GET", "/bar", jest.fn());
    expectation.addHandler("GET", "/baz/foo", jest.fn());

    expectation.addMiddleware("/", mware(jest.fn()));

    // First branch
    expectation.addHandler("POST", "/baz/foo", jest.fn());
    expectation.addHandler("PUT", "/baz/foo", jest.fn());

    // Second branch
    expectation.addHandler("POST", "/foo/:bar", jest.fn());
    expectation.addHandler("PUT", "/foo/:baz", jest.fn());

    expectation.addMiddleware("/foo", mware(jest.fn()));
    expectation.addMiddleware("/foo/foo", mware(jest.fn()));

    // First branch of second branch
    expectation.addHandler("GET", "/foo/:blah/foo/blah", jest.fn());
    expectation.addHandler("PUT", "/foo/:blah/baz", jest.fn());

    expectation.addMiddleware("/foo/:blah", mware(jest.fn()));
    expectation.addMiddleware("/foo/:blah/metrics", mware(jest.fn()));

    const flattened = topLevel.flatten();
    // expect(flattened).toEqual(expectation);
    compareTrees(expectation, flattened);
    // This should work both ways round, i.e. we want to make sure that
    // no extra routes have appeared in the flattened table!!
    compareTrees(flattened, expectation);
  });
});

function mware(fn: any): MiddlewareObj {
  return {
    handle: fn,
    type: MiddlewareT.MIDDLEWARE_FN,
  };
}

function compareTrees(expectation: RouteTree, actual: RouteTree) {
  for (const [route, methods] of expectation.handlers) {
    for (const [method, _handler] of methods) {
      expect(actual.handlers.get(route)!.get(method)).toEqual(
        expect.any(Function)
      );
    }
  }
  for (const [route, middlewares] of expectation.middlewares) {
    for (const [index, _middleware] of middlewares.entries()) {
      expect(actual.middlewares.get(route)![index]).toEqual({
        handle: expect.any(Function),
        type: MiddlewareT.MIDDLEWARE_FN,
      });
    }
  }
}
