package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	sv "github.com/foldsh/fold/runtime/supervisor"
)

const port = ":12345"

func TestServeHTTP(t *testing.T) {
	logger := logging.NewTestLogger()
	methodNotAllowed := map[string]interface{}{"title": "Method not allowed"}
	resourceNotFound := map[string]interface{}{"title": "Resource not found"}
	cases := []struct {
		name     string
		manifest *manifest.Manifest
		requests []testReq
	}{
		{
			"PlainPath",
			mkmanifest(mkroute("GET", "get", "/get")),
			[]testReq{
				{
					"GET",
					"/get",
					200,
					map[string]interface{}{
						"handler": "get",
						"path":    "/get",
					},
				},
				{"POST", "/get", 405, methodNotAllowed},
				{"PUT", "/get", 405, methodNotAllowed},
				{"GET", "/foo", 404, resourceNotFound},
			},
		},
		{
			"PathParams",
			mkmanifest(mkroute("GET", "foo", "/foo/:bar/:baz")),
			[]testReq{
				{
					"GET",
					"/foo/a/b",
					200,
					map[string]interface{}{
						"handler": "foo",
						"path":    "/foo/:bar/:baz",
						"bar":     "a",
						"baz":     "b",
					},
				},
				{
					"GET",
					"/foo/y/z",
					200,
					map[string]interface{}{
						"handler": "foo",
						"path":    "/foo/:bar/:baz",
						"bar":     "y",
						"baz":     "z",
					},
				},
				{"GET", "/baz/y/z", 404, resourceNotFound},
				{"POST", "/foo/y/z", 405, methodNotAllowed},
			},
		},
		{
			"QueryParams",
			mkmanifest(mkroute("GET", "foo", "/foo")),
			[]testReq{
				{
					"GET",
					"/foo?bar=a&baz=b",
					200,
					map[string]interface{}{
						"handler": "foo",
						"path":    "/foo",
						"bar":     "a",
						"baz":     "b",
					},
				},
				{
					"GET",
					"/foo?bar=y&baz=z",
					200,
					map[string]interface{}{
						"handler": "foo",
						"path":    "/foo",
						"bar":     "y",
						"baz":     "z",
					},
				},
				{
					"GET",
					"/foo?bar=a&bar=b&bar=c",
					200,
					map[string]interface{}{
						"handler": "foo",
						"path":    "/foo",
						// It's an array of interface{} because it makes using
						// go-cmp with unmarhsaled json much easier.
						"bar": []interface{}{"a", "b", "c"},
					},
				},
				{"GET", "/foo/y/z", 404, resourceNotFound},
			},
		},
		{
			"MultipleRoutes",
			mkmanifest(
				mkroute("GET", "getFoo", "/foo"),
				mkroute("PUT", "putFoo", "/foo"),
				mkroute("GET", "getBar", "/bar"),
				mkroute("DELETE", "deleteBar", "/bar"),
			),
			[]testReq{
				{
					"GET",
					"/foo",
					200,
					map[string]interface{}{
						"handler": "getFoo",
						"path":    "/foo",
					},
				},
				{
					"PUT",
					"/foo",
					200,
					map[string]interface{}{
						"handler": "putFoo",
						"path":    "/foo",
					},
				},
				{"DELETE", "/foo", 405, methodNotAllowed},
				{
					"GET",
					"/bar",
					200,
					map[string]interface{}{
						"handler": "getBar",
						"path":    "/bar",
					},
				},
				{
					"DELETE",
					"/bar",
					200,
					map[string]interface{}{
						"handler": "deleteBar",
						"path":    "/bar",
					},
				},
				{"PUT", "/bar", 405, methodNotAllowed},
				{"GET", "/baz", 404, resourceNotFound},
			},
		},
		{
			"Complex",
			mkmanifest(
				mkroute("GET", "complex", "/complex/:foo/query/:bar"),
			),
			[]testReq{
				{
					"GET",
					"/complex/hello/query/world?baz=1&baz=2&baz=3&solo=param",
					200,
					map[string]interface{}{
						"handler": "complex",
						"path":    "/complex/:foo/query/:bar",
						"foo":     "hello",
						"bar":     "world",
						"baz":     []interface{}{"1", "2", "3"},
						"solo":    "param",
					},
				},
			},
		},
	}
	ms := &mockSupervisor{t}
	router := NewRouter(logger, ms)
	go func() {
		http.ListenAndServe(port, router)
	}()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router.Configure(tc.manifest)
			for _, r := range tc.requests {
				t.Run(fmt.Sprintf("%s:%s", r.method, r.path), func(t *testing.T) {
					status, body := req(logger, t, r.method, r.path)
					if status != r.expectedStatus {
						t.Errorf("expected status of %d but found %d", r.expectedStatus, status)
					}
					testutils.Diff(
						t,
						r.expectedBody,
						testutils.UnmarshalJSON(t, body),
						"Body did not match expectation",
					)
				})
			}
		})
	}
}

type testReq struct {
	method         string
	path           string
	expectedStatus int
	expectedBody   map[string]interface{}
}

type mockSupervisor struct {
	t *testing.T
}

func (ms *mockSupervisor) DoRequest(req *sv.Request) (*sv.Response, error) {
	// This doesn't do anything useful but by sending back this mix of
	// of data in the response it's very easy to verify that the request
	// made its way through here appropriately.
	body := make(map[string]interface{})
	body["handler"] = req.Handler
	body["path"] = req.Path
	for key, value := range req.PathParams {
		body[key] = value
	}
	for key, value := range req.QueryParams {
		if len(value) == 1 {
			body[key] = value[0]
		} else {
			body[key] = value
		}
	}
	return &sv.Response{Status: 200, Body: testutils.MarshalJSON(ms.t, body)}, nil
}

func mkmanifest(routes ...*manifest.Route) *manifest.Manifest {
	return &manifest.Manifest{
		Name:   "TEST",
		Routes: routes,
	}
}

func mkroute(method, handler string, pathSpec string) *manifest.Route {
	httpMethod, err := manifest.HttpMethodFromString(method)
	if err != nil {
		panic(err)
	}
	return &manifest.Route{
		HttpMethod: httpMethod,
		Handler:    handler,
		PathSpec:   pathSpec,
	}
}

func req(l logging.Logger, t *testing.T, method string, path string) (int, []byte) {
	client := http.Client{}
	request, err := http.NewRequest(method, fmt.Sprintf("http://localhost%s%s", port, path), nil)
	if method == "PUT" || method == "POST" {
		request.Header.Add("Content-Type", "application/json")
	}
	l.Debug("request", request)
	if err != nil {
		t.Fatalf("failed to create request: %v ", err)
	}
	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("failed to do request: %v ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v ", err)
	}
	return resp.StatusCode, body
}
