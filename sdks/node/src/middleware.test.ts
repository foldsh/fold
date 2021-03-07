import bodyParser from "body-parser";
import cookieParser from "cookie-parser";
import signature from "cookie-signature";
import timeout from "connect-timeout";
import promBundle from "express-prom-bundle";
import { fold } from "./service";
import run from "./__tests__/run";

describe("Middleware", () => {
  test("should be able register middleware", async () => {
    const svc = fold();
    svc.use((req, res, next) => {
      req.params.foo = "bar";
      res.set("Vary", "User-Agent");
      next();
    });
    svc.get("/", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    await run(svc)
      .get({ route: "/" })
      .expectStatus(200)
      .expectBody({ foo: "bar" })
      .done();
  });
  test("should be able to nest services", async () => {
    const svc = fold();
    const sub = fold();
    sub.use((req, res, next) => {
      req.params.foo = "bar";
      res.set("Vary", "User-Agent");
      next();
    });
    sub.get("/bar", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    svc.use("/foo", sub);
    await run(svc)
      .get({ route: "/foo/bar" })
      .expectStatus(200)
      .expectBody({ foo: "bar" })
      .done();
  });
  test("should be able to next service on the default path", async () => {
    const svc = fold();
    const sub = fold();
    sub.use((req, res, next) => {
      req.params.foo = "bar";
      res.set("Vary", "User-Agent");
      next();
    });
    sub.get("/bar", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    svc.use(sub);
    await run(svc)
      .get({ route: "/bar" })
      .expectStatus(200)
      .expectBody({ foo: "bar" })
      .done();
  });
  test("should be able to end a request early from middleware", async () => {
    const svc = fold();
    const sub = fold();
    sub.use((_req, res, _next) => {
      res.send("early exit");
    });
    sub.get("/bar", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    svc.use(sub);
    await run(svc).get({ route: "/bar" }).expectBody("early exit").done();
  });
  test("should be able to use middleware on a path without a handler", async () => {
    // It's important to note the expected behaviour of the runtime here.
    // Fold does not support 'undefined' routes. I.e. only requests that match an explicitly
    // specified route will be forwarded on to the application. This is a constraint,
    // but also makes life simpler for the developer. They don't need to worry about
    // requests that don't match a route getting lost in the router.
    // Dynamic matching of paths is used fairly commonly in the express community but
    // fold doesn't do that. If you want a middleware to match a route you must explicitly
    // register it on that route.
    const svc = fold();
    svc.use("/metrics", (req, res, _next) => {
      if (req.path === "/metrics") {
        res.status(200);
        res.send({ responseTime: 100 });
      }
    });
    await run(svc)
      .get({ route: "/metrics", url: "https://www.fold.sh/metrics" })
      .expectBody({ responseTime: 100 })
      .done();
  });
  test("should be able to use nested middleware on a path without a handler", async () => {
    const svc = fold();
    const sub = fold();
    sub.use("/metrics", (_req, res, _next) => {
      res.status(200);
      res.send({ responseTime: 100 });
    });
    svc.use("/users", sub);
    await run(svc)
      .get({
        route: "/users/metrics",
        url: "https://www.fold.sh/users/metrics",
      })
      .expectBody({ responseTime: 100 })
      .done();
  });
  test("middleware should only be invoked on handlers beneath after it", async () => {
    const svc = fold();
    svc.get("/baz", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    const sub = fold();
    sub.use((req, res, next) => {
      req.params.foo = "bar";
      res.set("Vary", "User-Agent");
      next();
    });
    sub.get("/bar", (req, res) => {
      res.status(200);
      res.json({ foo: req.params.foo });
    });
    svc.use("/foo", sub);
    await run(svc)
      .get({ route: "/baz" })
      .expectStatus(200)
      .expectBody({ foo: undefined })
      .done();
  });
  describe("3rd Party Middlewares", () => {
    describe("body-parser", () => {
      test("should parse a json body", async () => {
        const svc = fold();
        svc.use(bodyParser.json());
        svc.get("/", (req, res) => {
          res.status(200);
          res.send({ requestBody: req.body });
        });
        await run(svc)
          .get({
            route: "/",
            body: { hello: "world" },
            headers: { "Content-Type": ["application/json"] },
          })
          .expectStatus(200)
          .expectHeader("Content-Type", "application/json; charset=utf-8")
          .expectBody({ requestBody: { hello: "world" } })
          .done();
      });
    });
    describe("cookie-parser", () => {
      test("should parse plain cookies", async () => {
        const svc = fold();
        svc.use(cookieParser() as any);
        svc.get("/", (req, res) => {
          res.status(200);
          res.send({ cookies: req.cookies });
        });
        await run(svc)
          .get({
            route: "/",
            headers: { Cookie: ["foo=bar; bar=baz"] },
          })
          .expectStatus(200)
          .expectHeader("Content-Type", "application/json; charset=utf-8")
          .expectBody({ cookies: { foo: "bar", bar: "baz" } })
          .done();
      });
      test("should parse json cookies", async () => {
        const svc = fold();
        svc.use(cookieParser() as any);
        svc.get("/", (req, res) => {
          res.status(200);
          res.send({ cookies: req.cookies });
        });
        await run(svc)
          .get({
            route: "/",
            headers: { Cookie: ['foo=j:{"bar":"baz"}'] },
          })
          .expectStatus(200)
          .expectHeader("Content-Type", "application/json; charset=utf-8")
          .expectBody({ cookies: { foo: { bar: "baz" } } })
          .done();
      });
      test("should parse signed cookies", async () => {
        const secret = "secret";
        const signed = signature.sign("foobar", secret);
        const svc = fold();
        svc.use(cookieParser(secret) as any);
        svc.get("/", (req, res) => {
          res.status(200);
          res.send({ cookies: req.signedCookies });
        });
        await run(svc)
          .get({
            route: "/",
            headers: { Cookie: [`foo=s:${signed}`] },
          })
          .expectStatus(200)
          .expectHeader("Content-Type", "application/json; charset=utf-8")
          .expectBody({ cookies: { foo: "foobar" } })
          .done();
      });
    });
    describe("connect-timeout", () => {
      jest.useFakeTimers();
      test("should timeout out a hanging requet", async () => {
        const svc = fold();
        // Requests should timeout after 1 second.
        svc.use(timeout("1s") as any);
        svc.use("/", (req, _res, next) => {
          // This middleware should execute within the timeout
          if (!req.timedout) {
            next();
          }
        });
        svc.get("/", (req, _res) => {
          // Advance time to trigger the timeout.
          jest.advanceTimersByTime(1000);
          // If the timeout hadn't triggered this would block forever.
          let timedout = req.timedout;
          while (!timedout) {
            timedout = req.timedout;
            continue;
          }
        });
        await run(svc)
          .get({
            route: "/",
          })
          .expectStatus(500)
          .expectBody({
            title: "Internal server error",
            detail: {
              code: "ETIMEDOUT",
              message: "Response timeout",
              timeout: 1000,
            },
          })
          .done();
      });
    });
    describe("prom-bundle", () => {
      const PROM_RESPONSE = `# HELP http_request_duration_seconds duration histogram of http responses labeled with: status_code, method
# TYPE http_request_duration_seconds histogram

# HELP up 1 = up, 0 = not up
# TYPE up gauge
up 1
`;
      test("should be able to fetch application metrics", async () => {
        const svc = fold();
        svc.use("/metrics", promBundle({ includeMethod: true }) as any);
        await run(svc)
          .get({ route: "/metrics", url: "https://www.fold.sh/metrics" })
          .expectStatus(200)
          .expectHeader("Content-Type", "text/plain; charset=utf-8")
          .expectBody(PROM_RESPONSE)
          .done();
      });
    });
  });
});
