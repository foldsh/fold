package types

import (
	"bytes"
	"net/http"

	"github.com/foldsh/fold/manifest"
)

// This is just a wrapper around the protobuf definition in proto/ingress.proto
// It makes them easier to use and avoids exposing the generated code to the
// rest of the runtime package.
type Request struct {
	HTTPMethod    string
	Path          string
	RawQuery      string
	Fragment      string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Host          string
	RemoteAddr    string
	RequestURI    string
	ContentLength int64
	Body          []byte
	Headers       map[string][]string
	PathParams    map[string]string
	QueryParams   map[string][]string
	Route         string
}

func ReqFromHTTP(req *http.Request, route string, pathParams map[string]string) *Request {
	buf := new(bytes.Buffer)
	// No need to close as a request body is closed by the server.
	buf.ReadFrom(req.Body)
	return &Request{
		HTTPMethod:    req.Method,
		Path:          req.URL.Path,
		RawQuery:      req.URL.RawQuery,
		Fragment:      req.URL.Fragment,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		RequestURI:    req.RequestURI,
		ContentLength: req.ContentLength,
		Body:          buf.Bytes(),
		Headers:       req.Header,
		PathParams:    pathParams,
		QueryParams:   req.URL.Query(),
		Route:         route,
	}
}

func (req *Request) ToProto() (*manifest.FoldHTTPRequest, error) {
	httpMethod, err := manifest.HTTPMethodFromString(req.HTTPMethod)
	if err != nil {
		return nil, err
	}
	return &manifest.FoldHTTPRequest{
		HttpMethod: httpMethod,
		Path:       req.Path,
		RawQuery:   req.RawQuery,
		Fragment:   req.Fragment,
		HttpProto: &manifest.FoldHTTPProto{
			Proto: req.Proto,
			Major: int32(req.ProtoMajor),
			Minor: int32(req.ProtoMinor),
		},
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		RequestUri:    req.RequestURI,
		ContentLength: req.ContentLength,
		Body:          req.Body,
		Headers:       encodeMapRepeatedString(req.Headers),
		PathParams:    req.PathParams,
		QueryParams:   encodeMapRepeatedString(req.QueryParams),
		Route:         req.Route,
	}, nil
}

// This is just a wrapper around the protobuf definition in proto/ingress.proto
// It makes them easier to use and avoids exposing the generated code to the
// rest of the runtime package.
type Response struct {
	Status  int
	Body    []byte
	Headers map[string][]string
}

func ResFromProto(res *manifest.FoldHTTPResponse) *Response {
	return &Response{
		Status:  int(res.Status),
		Body:    res.Body,
		Headers: decodeMapRepeatedString(res.Headers),
	}
}

func encodeMapRepeatedString(m map[string][]string) map[string]*manifest.StringArray {
	result := map[string]*manifest.StringArray{}
	for key, value := range m {
		result[key] = &manifest.StringArray{Values: value}
	}
	return result
}

func decodeMapRepeatedString(m map[string]*manifest.StringArray) map[string][]string {
	result := map[string][]string{}
	for key, value := range m {
		result[key] = value.Values
	}
	return result
}
