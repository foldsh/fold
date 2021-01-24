package router

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	sv "github.com/foldsh/fold/runtime/supervisor"
)

type Router interface {
	http.Handler
	// While this method signature is the same as ServeHTTP, it is
	// intended to support a different use case altogether.
	DoRequest(http.ResponseWriter, *http.Request)
	Configure(*manifest.Manifest)
}

type RequestDoer interface {
	DoRequest(*sv.Request) (*sv.Response, error)
}

// Builds a router from a service manifest. While we could fetch the manfiest
// from the service, making it a parameter gives some more options about
// how and when we acquire one.
func NewRouter(logger logging.Logger, doer RequestDoer) Router {
	return &foldRouter{logger: logger, doer: doer}
}

type foldRouter struct {
	logger logging.Logger
	doer   RequestDoer
	router *httprouter.Router
}

// This just implements the http.Handler interface
func (fr *foldRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fr.router.ServeHTTP(w, r)
}

func (fr *foldRouter) DoRequest(w http.ResponseWriter, r *http.Request) {
	// We will have to implement an appropriate response writer over in
	// the lambda code which uses that interface to produce a response for
	// the lambda proxy integration:
	// {
	//   "isBase64Encoded": False,
	//   "statusCode": res.get("status"),
	//   "headers": {},
	//   "body": json.dumps(res),
	// }
	// Choosing to ignore the trailing slash redirect for now.
	handle, ps, _ := fr.router.Lookup(r.Method, r.URL.Path)
	handle(w, r, ps)
}

func (fr *foldRouter) Configure(m *manifest.Manifest) {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFound)
	router.MethodNotAllowed = http.HandlerFunc(notAllowed)
	for _, route := range m.Routes {
		router.Handle(
			route.HttpMethod.String(),
			route.PathSpec,
			fr.makeHandler(route),
		)
	}
	fr.router = router
}

func (fr *foldRouter) makeHandler(route *manifest.Route) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method == "PUT" || r.Method == "POST" {
			isJSON := false
			for _, c := range r.Header.Values("Content-Type") {
				if c == "application/json" {
					isJSON = true
					break
				}
			}
			if !isJSON {
				unsupportedMediaType(w, r)
				return
			}
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		req := sv.Request{
			HttpMethod:  r.Method,
			Handler:     route.Handler,
			Path:        route.PathSpec,
			Body:        buf.Bytes(),
			Headers:     r.Header,
			PathParams:  encodePathParams(ps),
			QueryParams: r.URL.Query(),
		}
		res, err := fr.doer.DoRequest(&req)
		if err != nil {
			// TODO write a 500 response
		}
		// Write the status code
		w.WriteHeader(int(res.Status))
		// Write the headers
		headers := w.Header()
		for key, values := range res.Headers {
			for _, value := range values {
				headers.Add(key, value)
			}
		}
		// Write the body
		body := []byte(res.Body)
		n, err := w.Write(body)
		if err != nil {
			// TODO write a 500 response
		}
		if n != len(body) {
			// TODO try to recover by writing more? or 500
		}
	}
}

func encodePathParams(params httprouter.Params) map[string]string {
	result := map[string]string{}
	for _, param := range params {
		result[param.Key] = param.Value
	}
	return result
}

func httpError(w http.ResponseWriter, code int, e string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, e)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusNotFound, `{"title":"Resource not found"}`)
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusMethodNotAllowed, `{"title":"Method not allowed"}`)
}

func unsupportedMediaType(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusUnsupportedMediaType, `{"title":"Content-Type must be application/json"}`)
}
