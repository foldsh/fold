package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/transport"
)

type RequestDoer interface {
	DoRequest(context.Context, *transport.Request) (*transport.Response, error)
}

// Builds a router from a service manifest. While we could fetch the manfiest
// from the service, making it a parameter gives some more options about
// how and when we acquire one.
func NewRouter(logger logging.Logger, doer RequestDoer) *Router {
	return &Router{logger: logger, doer: doer, router: newRouter()}
}

var HTTP_METHODS = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

func NewCatchAllRouter(logger logging.Logger, handler http.Handler) *Router {
	router := httprouter.New()
	for _, method := range HTTP_METHODS {
		router.Handler(method, "/*all", handler)
	}
	return &Router{logger: logger, router: router}
}

type Router struct {
	logger   logging.Logger
	doer     RequestDoer
	router   *httprouter.Router
	manifest *manifest.Manifest
}

// This just implements the http.Handler interface
func (fr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fr.router.ServeHTTP(w, r)
}

func (fr *Router) Configure(m *manifest.Manifest) {
	fr.manifest = m
	router := newRouter()
	// Register the default admin routes.
	router.Handle(
		"GET",
		"/_foldadmin/healthz",
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK"}`))
		},
	)
	router.Handle(
		"GET",
		"/_foldadmin/manifest",
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// If successful Write implicity sets the 200 response on the ResponseWriter
			err := manifest.WriteJSON(w, fr.manifest)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(`{"title":"Failed to marshal manifest to JSON"}`))
				return
			}
		},
	)
	// And now register all of the routes from the manifest.
	for _, route := range m.Routes {
		router.Handle(
			route.HttpMethod.String(),
			route.Route,
			fr.makeHandler(route),
		)
	}
	fr.router = router
}

func newRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFound)
	router.MethodNotAllowed = http.HandlerFunc(notAllowed)
	return router
}

func (fr *Router) makeHandler(route *manifest.Route) httprouter.Handle {
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
		req := transport.ReqFromHTTP(r, route.Route, encodePathParams(ps))
		res, err := fr.doer.DoRequest(r.Context(), req)
		if err != nil {
			httpError(
				w,
				500,
				fmt.Sprintf(`{"title": "Runtime error", "detail": "%v"}`, err),
			)
			return
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
			httpError(w, 500, fmt.Sprintf(`{"title": "Runtime error", "detail": "%v"}`, err))
			return
		}
		if n != len(body) {
			httpError(w, 500, `{"title": "Runtime error", "detail": "Failed to read entire body."}`)
			return
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
	httpError(
		w,
		http.StatusUnsupportedMediaType,
		`{"title":"Content-Type must be application/json"}`,
	)
}
