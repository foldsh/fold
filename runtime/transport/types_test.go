package transport_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/transport"
)

func TestRequestToProto(t *testing.T) {
	req := &transport.Request{
		HTTPMethod:    "GET",
		Path:          "/foo/bar",
		RawQuery:      "foo=bar",
		Fragment:      "foo/bar",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Host:          "fold.sh",
		RemoteAddr:    "remote.com",
		RequestURI:    "/foo/bar?foo=bar#foo/bar",
		ContentLength: 1,
		Body:          []byte("a"),
		Headers:       map[string][]string{"Vary": []string{"User-Agent"}},
		PathParams:    map[string]string{"baz": "bar"},
		QueryParams:   map[string][]string{"foo": []string{"bar"}},
		Route:         "/foo/:baz",
	}
	expectation := &manifest.FoldHTTPRequest{
		HttpMethod: manifest.FoldHTTPMethod_GET,
		Path:       "/foo/bar",
		RawQuery:   "foo=bar",
		Fragment:   "foo/bar",
		HttpProto: &manifest.FoldHTTPProto{
			Proto: "HTTP/1.1",
			Major: 1,
			Minor: 1,
		},
		Host:          "fold.sh",
		RemoteAddr:    "remote.com",
		RequestUri:    "/foo/bar?foo=bar#foo/bar",
		ContentLength: 1,
		Body:          []byte("a"),
		Headers: map[string]*manifest.StringArray{
			"Vary": &manifest.StringArray{Values: []string{"User-Agent"}},
		},
		PathParams: map[string]string{"baz": "bar"},
		QueryParams: map[string]*manifest.StringArray{
			"foo": &manifest.StringArray{Values: []string{"bar"}},
		},
		Route: "/foo/:baz",
	}
	result, err := req.ToProto()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	testutils.DiffFoldHTTPRequest(t, expectation, result)
}

func TestRequestFromHTTP(t *testing.T) {
	url := "http://www.fold.sh/foo/bar?foo=bar#foo/bar"
	req, err := http.NewRequest(
		"GET", url,
		strings.NewReader("foo"),
	)
	req.Header["Vary"] = []string{"User-Agent"}
	if err != nil {
		t.Errorf("%+v", err)
	}
	req.RequestURI = url
	expectation := &transport.Request{
		HTTPMethod:    "GET",
		Path:          "/foo/bar",
		RawQuery:      "foo=bar",
		Fragment:      "foo/bar",
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Host:          "www.fold.sh",
		RemoteAddr:    "",
		RequestURI:    url,
		ContentLength: 3,
		Body:          []byte("foo"),
		Headers: map[string][]string{
			"Vary": []string{"User-Agent"},
		},
		PathParams:  map[string]string{"baz": "bar"},
		QueryParams: map[string][]string{"foo": []string{"bar"}},
		Route:       "/foo/:baz",
	}
	result := transport.ReqFromHTTP(req, "/foo/:baz", map[string]string{"baz": "bar"})
	testutils.Diff(t, expectation, result, "Response does not match expectation")
}

func TestResponseFromProto(t *testing.T) {
	res := &manifest.FoldHTTPResponse{
		Status: int32(200),
		Body:   []byte("foo"),
		Headers: map[string]*manifest.StringArray{
			"Content-Length": &manifest.StringArray{Values: []string{"3"}},
			"Content-Type": &manifest.StringArray{
				Values: []string{"application/json; charset=utf-8"},
			},
		},
	}
	expectation := &transport.Response{
		Status: 200,
		Body:   []byte("foo"),
		Headers: map[string][]string{
			"Content-Length": []string{"3"},
			"Content-Type":   []string{"application/json; charset=utf-8"},
		},
	}
	result := transport.ResFromProto(res)
	testutils.Diff(t, expectation, result, "Response does not match expectation")
}
