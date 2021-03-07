import { fold } from "./service";
import run from "./__tests__/run";

describe("grpc.getManifest", () => {
  test("should return a simple manifest", async () => {
    const svc = fold({ major: 1, minor: 2, patch: 3 });
    svc.get("/", (_req, res) => {
      res.sendStatus(200);
    });
    await run(svc)
      .manifest()
      .expectManifest({
        version: { major: 1, minor: 2, patch: 3 },
        routes: [{ method: "GET", handler: "GET /", route: "/" }],
      })
      .done();
  });
  test("should return a manifest from nested services", async () => {
    const svc = fold({ major: 1, minor: 2, patch: 3 });
    svc.get("/", (_req, res) => {
      res.sendStatus(200);
    });
    const foo = fold();
    const bar = fold();
    const baz = fold();
    foo.get("/get", (_req, _res) => {});
    foo.put("/put", (_req, _res) => {});
    svc.use("/foo", foo);

    bar.delete("/delete", (_req, _res) => {});
    bar.post("/post", (_req, _res) => {});

    baz.patch("/patch", (_req, _res) => {});
    baz.get("/get", (_req, _res) => {});

    baz.use("/bar", bar);
    svc.use("/baz", baz);
    await run(svc)
      .manifest()
      .expectManifest({
        version: { major: 1, minor: 2, patch: 3 },
        routes: [
          { method: "GET", handler: "GET /", route: "/" },
          { method: "GET", handler: "GET /foo/get", route: "/foo/get" },
          { method: "PUT", handler: "PUT /foo/put", route: "/foo/put" },
          {
            method: "PATCH",
            handler: "PATCH /baz/patch",
            route: "/baz/patch",
          },
          { method: "GET", handler: "GET /baz/get", route: "/baz/get" },
          {
            method: "DELETE",
            handler: "DELETE /baz/bar/delete",
            route: "/baz/bar/delete",
          },
          {
            method: "POST",
            handler: "POST /baz/bar/post",
            route: "/baz/bar/post",
          },
        ],
      })
      .done();
  });
});
