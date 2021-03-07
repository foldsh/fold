package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/types"
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
			"Should match simple routes.",
			mkmanifest(mkroute("GET", "/get")),
			[]testReq{
				{
					method:         "GET",
					path:           "/get",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/get",
						"path":  "/get",
					},
				},
				{
					method:         "POST",
					path:           "/get",
					expectedStatus: 405,
					expectedBody:   methodNotAllowed,
				},
				{
					method:         "PUT",
					path:           "/get",
					expectedStatus: 405,
					expectedBody:   methodNotAllowed,
				},
				{
					method:         "GET",
					path:           "/foo",
					expectedStatus: 404,
					expectedBody:   resourceNotFound,
				},
			},
		},
		{
			"Should match routes with parameters",
			mkmanifest(mkroute("GET", "/foo/:bar/:baz")),
			[]testReq{
				{
					method:         "GET",
					path:           "/foo/a/b",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"path":  "/foo/a/b",
						"route": "/foo/:bar/:baz",
						"bar":   "a",
						"baz":   "b",
					},
				},
				{
					method:         "GET",
					path:           "/foo/y/z",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"path":  "/foo/y/z",
						"route": "/foo/:bar/:baz",
						"bar":   "y",
						"baz":   "z",
					},
				},
				{
					method:         "GET",
					path:           "/baz/y/z",
					expectedStatus: 404,
					expectedBody:   resourceNotFound,
				},
				{
					method:         "POST",
					path:           "/foo/y/z",
					expectedStatus: 405,
					expectedBody:   methodNotAllowed,
				},
			},
		},
		{
			"Should parse query params",
			mkmanifest(mkroute("GET", "/foo")),
			[]testReq{
				{
					method:         "GET",
					path:           "/foo?bar=a&baz=b",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/foo",
						"path":  "/foo",
						"bar":   "a",
						"baz":   "b",
					},
				},
				{
					method:         "GET",
					path:           "/foo?bar=y&baz=z",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/foo",
						"path":  "/foo",
						"bar":   "y",
						"baz":   "z",
					},
				},
				{
					method:         "GET",
					path:           "/foo?bar=a&bar=b&bar=c",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/foo",
						"path":  "/foo",
						// It's an array of interface{} because it makes using
						// go-cmp with unmarhsaled json much easier.
						"bar": []interface{}{"a", "b", "c"},
					},
				},
				{
					method:         "GET",
					path:           "/foo/y/z",
					expectedStatus: 404,
					expectedBody:   resourceNotFound,
				},
			},
		},
		{
			"Should allow multiple routes to be set",
			mkmanifest(
				mkroute("GET", "/foo"),
				mkroute("PUT", "/foo"),
				mkroute("GET", "/bar"),
				mkroute("DELETE", "/bar"),
			),
			[]testReq{
				{
					method:         "GET",
					path:           "/foo",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/foo",
						"path":  "/foo",
					},
				},
				{
					method:         "PUT",
					path:           "/foo",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/foo",
						"path":  "/foo",
					},
				},
				{
					method:         "DELETE",
					path:           "/foo",
					expectedStatus: 405,
					expectedBody:   methodNotAllowed,
				},
				{
					method:         "GET",
					path:           "/bar",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/bar",
						"path":  "/bar",
					},
				},
				{
					method:         "DELETE",
					path:           "/bar",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"route": "/bar",
						"path":  "/bar",
					},
				},
				{
					method:         "PUT",
					path:           "/bar",
					expectedStatus: 405,
					expectedBody:   methodNotAllowed,
				},
				{
					method:         "GET",
					path:           "/baz",
					expectedStatus: 404,
					expectedBody:   resourceNotFound,
				},
			},
		},
		{
			"Should process a complex route correctly",
			mkmanifest(
				mkroute("GET", "/complex/:foo/query/:bar"),
			),
			[]testReq{
				{
					method:         "GET",
					path:           "/complex/hello/query/world?baz=1&baz=2&baz=3&solo=param",
					expectedStatus: 200,
					expectedBody: map[string]interface{}{
						"path":  "/complex/hello/query/world",
						"route": "/complex/:foo/query/:bar",
						"foo":   "hello",
						"bar":   "world",
						"baz":   []interface{}{"1", "2", "3"},
						"solo":  "param",
					},
				},
			},
		},
		{
			"Should pass on a request body",
			mkmanifest(mkroute("POST", "/post")),
			[]testReq{
				{
					method:         "POST",
					path:           "/post",
					expectedStatus: 200,
					body: map[string]interface{}{
						"foo": "bar",
					},
					expectedBody: map[string]interface{}{
						"path":  "/post",
						"route": "/post",
						"body": map[string]interface{}{
							"foo": "bar",
						},
					},
				},
			},
		},
	}
	ms := &mockRequestDoer{t}
	router := NewRouter(logger, ms)
	go func() {
		http.ListenAndServe(port, router)
	}()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router.Configure(tc.manifest)
			// Test healthz
			status, body := req(logger, t, "GET", "/_foldadmin/healthz", nil)
			if status != 200 {
				t.Errorf("Expected 200 response from healthz but found %d", status)
			}
			testutils.Diff(
				t,
				`{"status":"OK"}`,
				string(body),
				"/_foldadmin/healthz body did not match expectation",
			)
			// Test the manifest
			expectation := &bytes.Buffer{}
			manifest.WriteJSON(expectation, tc.manifest)
			status, actual := req(logger, t, "GET", "/_foldadmin/manifest", nil)
			if status != 200 {
				t.Errorf("Expected 200 response from manifest but found %d", status)
			}
			testutils.Diff(
				t,
				expectation.Bytes(),
				actual,
				"/_foldadmin/manifest did not return the expected manifest",
			)
			for _, r := range tc.requests {
				t.Run(fmt.Sprintf("%s:%s", r.method, r.path), func(t *testing.T) {
					status, body := req(logger, t, r.method, r.path, r.body)
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
	body           map[string]interface{}
	expectedStatus int
	expectedBody   map[string]interface{}
}

type mockRequestDoer struct {
	t *testing.T
}

func (ms *mockRequestDoer) DoRequest(req *types.Request) (*types.Response, error) {
	// First up we'll run some generic assertions about the request being
	// formed.
	if req.Proto != "HTTP/1.1" {
		ms.t.Errorf("Expected HTTP1.1 but found %s", req.Proto)
	}
	if req.ProtoMajor != 1 {
		ms.t.Errorf("Expected proto major of 1 but found %d", req.ProtoMajor)
	}
	if req.ProtoMinor != 1 {
		ms.t.Errorf("Expected proto minor of 1 but found %d", req.ProtoMinor)
	}
	if int(req.ContentLength) != len(req.Body) {
		ms.t.Errorf("ContentLength should equal body length")
	}
	if req.Host != "localhost:12345" {
		ms.t.Errorf("Expected host of localhost:12345 but found %s", req.Host)
	}
	// This doesn't do anything useful but by sending back this mix of
	// of data in the response it's very easy to verify that the request
	// made its way through here appropriately.
	body := make(map[string]interface{})
	body["route"] = req.Route
	body["path"] = req.Path
	if len(req.Body) != 0 {
		body["body"] = testutils.UnmarshalJSON(ms.t, req.Body)
	}
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
	return &types.Response{Status: 200, Body: testutils.MarshalJSON(ms.t, body)}, nil
}

func mkmanifest(routes ...*manifest.Route) *manifest.Manifest {
	return &manifest.Manifest{
		Name:   "TEST",
		Routes: routes,
	}
}

func mkroute(method, route string) *manifest.Route {
	httpMethod, err := manifest.HTTPMethodFromString(method)
	if err != nil {
		panic(err)
	}
	return &manifest.Route{
		HttpMethod: httpMethod,
		Route:      route,
	}
}

func req(
	l logging.Logger,
	t *testing.T,
	method string,
	path string,
	body map[string]interface{},
) (int, []byte) {
	client := http.Client{}
	var bodyReader io.Reader
	if body != nil {
		jsonBody := testutils.MarshalJSON(t, body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}
	request, err := http.NewRequest(
		method,
		fmt.Sprintf("http://localhost%s%s", port, path),
		bodyReader,
	)
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
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v ", err)
	}
	return resp.StatusCode, resBody
}
