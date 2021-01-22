package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	sv "github.com/foldsh/fold/runtime/supervisor"
)

const port = ":12345"

func TestServeHTTP(t *testing.T) {
	logger := logging.NewTestLogger()
	methodNotAllowed := "Method Not Allowed\n"
	pageNotFound := "404 page not found\n"
	cases := []struct {
		name     string
		manifest *manifest.Manifest
		requests []testReq
	}{
		{
			"PlainPath",
			mm(r(manifest.HttpMethod_GET, "get", "/get")),
			[]testReq{
				{"GET", "/get", 200, "handler:get path:/get "},
				{"POST", "/get", 405, methodNotAllowed},
				{"PUT", "/get", 405, methodNotAllowed},
				{"GET", "/foo", 404, pageNotFound},
			},
		},
		{
			"PathParams",
			mm(r(manifest.HttpMethod_GET, "foo", "/foo/:bar/:baz")),
			[]testReq{
				{"GET", "/foo/a/b", 200, "handler:foo path:/foo/:bar/:baz bar:a baz:b "},
				{"GET", "/foo/y/z", 200, "handler:foo path:/foo/:bar/:baz bar:y baz:z "},
				{"GET", "/baz/y/z", 404, pageNotFound},
				{"POST", "/foo/y/z", 405, methodNotAllowed}},
		},
	}

	router := NewRouter(logger, &mockSupervisor{})
	go func() {
		http.ListenAndServe(port, router)
	}()
	// Lame but just lets the server start
	// time.Sleep(10 * time.Millisecond)
	for _, tc := range cases {
		router.Configure(tc.manifest)
		t.Run(tc.name, func(t *testing.T) {
			for _, r := range tc.requests {
				status, body := req(logger, t, r.method, r.path)
				if status != r.status {
					t.Fatalf("expected status of %d but found %d", r.status, status)
				}
				if body != r.body {
					t.Fatalf("expected body of %s but found %s", r.body, body)
				}
			}
		})
	}
}

type testReq struct {
	method string
	path   string
	status int
	body   string
}

type mockSupervisor struct{}

func (ms *mockSupervisor) DoRequest(req *sv.Request) (*sv.Response, error) {
	var body strings.Builder
	fmt.Fprintf(&body, "handler:%s ", req.Handler)
	fmt.Fprintf(&body, "path:%s ", req.Path)
	for key, value := range req.PathParams {
		fmt.Fprintf(&body, "%s:%s ", key, value)
	}
	// Sticking the path back in the body with the request params and
	// setting the status code gives us evidence that the request has
	// made it through to here and has been evaluated
	return &sv.Response{Status: 200, Body: []byte(body.String())}, nil
}

func mm(routes ...*manifest.Route) *manifest.Manifest {
	return &manifest.Manifest{
		Name:   "TEST",
		Routes: routes,
	}
}

func r(method manifest.HttpMethod, handler string, pathSpec string) *manifest.Route {
	return &manifest.Route{HttpMethod: method, Handler: handler, PathSpec: pathSpec}
}

func req(l logging.Logger, t *testing.T, method string, path string) (int, string) {
	client := http.Client{}
	request, err := http.NewRequest(method, fmt.Sprintf("http://localhost%s%s", port, path), nil)
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
	return resp.StatusCode, string(body)
}
