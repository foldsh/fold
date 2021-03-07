package fold

import (
	"fmt"
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

type Request struct {
	HTTPMethod  string
	Body        map[string]interface{}
	Headers     map[string][]string
	PathParams  map[string]string
	QueryParams map[string][]string
	Route       string
}

type Response struct {
	StatusCode int
	Body       map[string]interface{}
	Headers    map[string][]string
}

type Handler func(*Request, *Response)

type Service interface {
	Start()
	Version(major, minor, patch int)
	Get(string, Handler)
	Put(string, Handler)
	Post(string, Handler)
	Delete(string, Handler)
	Logger() logging.Logger
}

func NewService() Service {
	stage := os.Getenv("FOLD_STAGE")
	name := os.Getenv("FOLD_SERVICE_NAME")
	var (
		logger logging.Logger
		err    error
	)
	switch stage {
	case "PRODUCTION":
		logger, err = logging.NewLogger(logging.Info, true)
	case "LOCAL":
		logger, err = logging.NewLogger(logging.Debug, false)
	default:
		logger, err = logging.NewLogger(logging.Debug, true)
	}
	if err != nil {
		panic(fmt.Sprintf("failed to start fold logger: %v", err))
	}
	s := &service{
		name:     name,
		handlers: make(map[string]map[string]Handler),
		logger:   logger,
		manifest: &manifest.Manifest{Name: name},
	}
	grpcServer := &grpcServer{service: s, logger: logger}
	s.server = grpcServer
	return s
}

type service struct {
	name     string
	server   *grpcServer
	manifest *manifest.Manifest
	handlers map[string]map[string]Handler
	logger   logging.Logger
}

func (s *service) Start() {
	s.logger.Infof("Starting fold service %v", s.name)
	s.server.start()
}

func (s *service) Version(major, minor, patch int) {
	s.manifest.Version = &manifest.Version{
		Major: int32(major),
		Minor: int32(minor),
		Patch: int32(patch),
	}
}

func (s *service) Get(route string, handler Handler) {
	s.registerHandler("GET", route, handler)
}

func (s *service) Head(route string, handler Handler) {
	s.registerHandler("HEAD", route, handler)
}

func (s *service) Post(route string, handler Handler) {
	s.registerHandler("POST", route, handler)
}

func (s *service) Put(route string, handler Handler) {
	s.registerHandler("PUT", route, handler)
}

func (s *service) Delete(route string, handler Handler) {
	s.registerHandler("DELETE", route, handler)
}

func (s *service) Connect(route string, handler Handler) {
	s.registerHandler("CONNECT", route, handler)
}

func (s *service) Options(route string, handler Handler) {
	s.registerHandler("OPTIONS", route, handler)
}

func (s *service) Trace(route string, handler Handler) {
	s.registerHandler("TRACE", route, handler)
}

func (s *service) Patch(route string, handler Handler) {
	s.registerHandler("PATCH", route, handler)
}

func (s *service) Logger() logging.Logger {
	return s.logger
}

func (s *service) registerHandler(method, route string, handler Handler) {
	// We can safely ignore the error because we control which strings it is possible to pass in .
	httpMethod, _ := manifest.HTTPMethodFromString(method)
	s.manifest.Routes = append(s.manifest.Routes, &manifest.Route{
		HttpMethod: httpMethod,
		Route:      route,
	})
	if _, exists := s.handlers[route]; !exists {
		s.handlers[route] = make(map[string]Handler)
	}
	s.handlers[route][method] = handler
}

func (s *service) doRequest(req *Request, res *Response) {
	if methods, exists := s.handlers[req.Route]; exists {
		if handler, exists := methods[req.HTTPMethod]; exists {
			handler(req, res)
			return
		}
	}
	// This shouldn't every happen as the runtime should only send handler ids
	// specified by the manifest.
	res.StatusCode = 500
	res.Body = map[string]interface{}{"title": fmt.Sprintf("Handler %s does not exist", req.Route)}
}
