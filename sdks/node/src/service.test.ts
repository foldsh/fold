import bodyParser from "body-parser";
import { fold } from "./service";
import run from "./__tests__/run";

describe("Service", () => {
  describe("svc.get", () => {
    test("should handle get requests", async () => {
      const svc = fold();
      svc.get("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .get({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
    test("should return json response bodies", async () => {
      const svc = fold();
      svc.get("/foo", (_req, res) => {
        res.status(200);
        res.json({ foo: "bar" });
      });
      await run(svc)
        .get({ route: "/foo" })
        .expectStatus(200)
        .expectBody({ foo: "bar" })
        .expectHeader("Content-Type", "application/json; charset=utf-8")
        .done();
    });
    test("should pass path params appropriately", async () => {
      const svc = fold();
      svc.get("/foo/:bar", (req, res) => {
        res.status(200);
        res.json({ foo: req.params.bar });
      });
      await run(svc)
        .get({ route: "/foo/:bar", params: { bar: "fold" } })
        .expectStatus(200)
        .expectBody({ foo: "fold" })
        .expectHeader("Content-Type", "application/json; charset=utf-8")
        .done();
    });
    test("should pass query params appropriately", async () => {
      const svc = fold();
      svc.get("/foo", (req, res) => {
        res.status(200);
        res.json({
          singleParamsAreAString: req.query.bar,
          multipleParamsAreAnArray: req.query.multiple,
        });
      });
      await run(svc)
        .get({
          route: "/foo",
          query: { bar: ["fold"], multiple: ["hello", "world"] },
        })
        .expectStatus(200)
        .expectBody({
          singleParamsAreAString: "fold",
          multipleParamsAreAnArray: ["hello", "world"],
        })
        .expectHeader("Content-Type", "application/json; charset=utf-8")
        .done();
    });
    test("should be able to set multiple routes", async () => {
      const svc = fold();
      svc.get("/foo", (_req, res) => {
        res.sendStatus(204);
      });
      svc.get("/bar", (_req, res) => {
        res.sendStatus(202);
      });
      await run(svc)
        .get({ route: "/foo" })
        .expectStatus(204)
        .expectBody("No Content")
        .done();
      await run(svc)
        .get({ route: "/bar" })
        .expectStatus(202)
        .expectBody("Accepted")
        .done();
    });
  });
  describe("svc.head", () => {
    test("should handle head requests", async () => {
      const svc = fold();
      svc.head("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .head({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.post", () => {
    test("should handle post requests", async () => {
      const svc = fold();
      svc.post("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .post({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
    test("should receive a json request body", async () => {
      const svc = fold();
      svc.use(bodyParser.json());
      svc.post("/", (req, res) => {
        res.status(200);
        res.json({ body: req.body });
      });
      const body = { foo: "bar" };
      await run(svc)
        .post({
          route: "/",
          body: body,
          headers: { "Content-Type": ["application/json"] },
        })
        .expectStatus(200)
        .expectBody({ body: body })
        .expectHeader("Content-Type", "application/json; charset=utf-8")
        .done();
    });
  });
  describe("svc.put", () => {
    test("should handle put requests", async () => {
      const svc = fold();
      svc.put("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .put({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.delete", () => {
    test("should handle delete requests", async () => {
      const svc = fold();
      svc.delete("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .delete({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.connect", () => {
    test("should handle connect requests", async () => {
      const svc = fold();
      svc.connect("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .connect({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.options", () => {
    test("should handle options requests", async () => {
      const svc = fold();
      svc.options("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .options({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.trace", () => {
    test("should handle trace requests", async () => {
      const svc = fold();
      svc.trace("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .trace({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.patch", () => {
    test("should handle patch requests", async () => {
      const svc = fold();
      svc.patch("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc)
        .patch({ route: "/" })
        .expectStatus(200)
        .expectBody("OK")
        .expectHeader("Content-Type", "text/plain; charset=utf-8")
        .done();
    });
  });
  describe("svc.all", () => {
    test("should respond to get request", async () => {
      const svc = fold();
      svc.all("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc).get({ route: "/" }).expectStatus(200).done();
    });
    test("should respond to put request", async () => {
      const svc = fold();
      svc.all("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc).put({ route: "/" }).expectStatus(200).done();
    });
    test("should respond to post request", async () => {
      const svc = fold();
      svc.all("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc).post({ route: "/" }).expectStatus(200).done();
    });
    test("should respond to delete request", async () => {
      const svc = fold();
      svc.all("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc).delete({ route: "/" }).expectStatus(200).done();
    });
    test("should respond to patch request", async () => {
      const svc = fold();
      svc.all("/", (_req, res) => {
        res.sendStatus(200);
      });
      await run(svc).patch({ route: "/" }).expectStatus(200).done();
    });
  });
  describe("svc.use", () => {
    test("should not accept just a path", () => {
      const svc = fold();
      expect(() => svc.use("/")).toThrowError(TypeError);
    });
    test("should not accept a path as the second parameter", () => {
      const svc = fold();
      // In ts the compiler complains without the cast to any but I am
      // including the test case for people writing in js.
      expect(() => svc.use((_req, _res, _next) => {}, "/" as any)).toThrowError(
        TypeError
      );
    });
    test("should be able to nest services with use", async () => {
      const foo = fold();
      foo.get("/", (_req, res) => {
        res.send("hello from foo");
      });
      const bar = fold();
      bar.use(foo);
      await run(bar).get({ route: "/" }).expectBody("hello from foo").done();
    });
    test("should be able to mount services on a path with use", async () => {
      const foo = fold();
      foo.get("/", (_req, res) => {
        res.send("hello from foo");
      });
      const bar = fold();
      bar.use("/foo", foo);
      await run(bar).get({ route: "/foo" }).expectBody("hello from foo").done();
    });
    test("should be able to mount services with many layers of nesting", async () => {
      const foo = fold();
      foo.get("/", (_req, res) => {
        res.send("hello from foo");
      });
      const bar = fold();
      const baz = fold();
      const root = fold();
      bar.use("/foo", foo);
      baz.use("/bar", bar);
      root.use("/baz", baz);
      await run(root)
        .get({ route: "/baz/bar/foo" })
        .expectBody("hello from foo")
        .done();
    });
    test("should be able to pass params through many services", async () => {
      const foo = fold();
      foo.get("/", (req, res) => {
        res.send(`${req.params.baz} from ${req.params.foo}`);
      });
      const bar = fold();
      const baz = fold();
      const root = fold();
      bar.use("/:foo", foo);
      baz.use("/bar", bar);
      root.use("/:baz", baz);
      await run(root)
        .get({
          route: "/:baz/bar/:foo",
          params: { baz: "hello", foo: "foo" },
        })
        .expectBody("hello from foo")
        .done();
    });
    test("should run middlewares registered on the root service", async () => {
      const svc = fold();
      svc.use((_req, res, next) => {
        res.status(301);
        next();
      });
      svc.get("/", (_req, res) => {
        res.send("redirect");
      });
      await run(svc)
        .get({ route: "/" })
        .expectStatus(301)
        .expectBody("redirect")
        .done();
    });
    test("should run multiple middlewares registered on the root service", async () => {
      const svc = fold();
      svc.use((_req, res, next) => {
        res.status(301);
        res.vary("Content-Type");
        next();
      });
      svc.use((_req, res, next) => {
        res.vary("User-Agent");
        next();
      });
      svc.get("/", (_req, res) => {
        res.send("redirect");
      });
      await run(svc)
        .get({ route: "/" })
        .expectStatus(301)
        .expectBody("redirect")
        .expectHeader("Vary", "Content-Type, User-Agent")
        .done();
    });
    test("should short circuit the request without a call to next", async () => {
      const svc = fold();
      svc.use((_req, res, _next) => {
        res.sendStatus(401);
      });
      svc.get("/", (_req, res) => {
        res.vary("User-Agent");
        res.send("handling");
      });
      await run(svc)
        .get({ route: "/" })
        .expectStatus(401)
        .expectBody("Unauthorized")
        .expectHeader("Vary", undefined)
        .done();
    });
  });
  describe("svc.settings.get", () => {});
  describe("svc.settings.set", () => {});
  describe("svc.settings.enable", () => {});
  describe("svc.settings.enabled", () => {});
  describe("svc.settings.disable", () => {});
});
