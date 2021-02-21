package fold

import (
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

type Request struct {
	HttpMethod  string
	Handler     string
	Path        string
	Body        map[string]interface{}
	Headers     map[string][]string
	PathParams  map[string]string
	QueryParams map[string][]string
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
		handlers: make(map[string]Handler),
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
	handlers map[string]Handler
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

func (s *service) Get(path string, handler Handler) {
	s.registerHandler("GET", path, handler)
}

func (s *service) Head(path string, handler Handler) {
	s.registerHandler("HEAD", path, handler)
}

func (s *service) Post(path string, handler Handler) {
	s.registerHandler("POST", path, handler)
}

func (s *service) Put(path string, handler Handler) {
	s.registerHandler("PUT", path, handler)
}

func (s *service) Delete(path string, handler Handler) {
	s.registerHandler("DELETE", path, handler)
}

func (s *service) Connect(path string, handler Handler) {
	s.registerHandler("CONNECT", path, handler)
}

func (s *service) Options(path string, handler Handler) {
	s.registerHandler("OPTIONS", path, handler)
}

func (s *service) Trace(path string, handler Handler) {
	s.registerHandler("TRACE", path, handler)
}

func (s *service) Patch(path string, handler Handler) {
	s.registerHandler("PATCH", path, handler)
}

func (s *service) Logger() logging.Logger {
	return s.logger
}

func (s *service) registerHandler(method, path string, handler Handler) {
	handlerId := uuid.NewV4().String()
	// We can safely ignore the error because we control which strings it is possible to pass in .
	httpMethod, _ := manifest.HttpMethodFromString(method)
	s.manifest.Routes = append(s.manifest.Routes, &manifest.Route{
		HttpMethod: httpMethod,
		PathSpec:   path,
		Handler:    handlerId,
	})
	s.handlers[handlerId] = handler
}

func (s *service) doRequest(req *Request, res *Response) {
	if handler, exists := s.handlers[req.Handler]; exists {
		handler(req, res)
	} else {
		// This shouldn't every happen as the runtime should only send handler ids
		// specified by the manifest.
		res.StatusCode = 500
		res.Body = map[string]interface{}{"title": fmt.Sprintf("Handler %s does not exist", req.Handler)}
	}
}
