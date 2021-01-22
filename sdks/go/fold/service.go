package fold

import (
	"github.com/satori/go.uuid"

	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/sdks/go/internal"
)

type Service struct {
	server   internal.GrpcServer
	manifest *manifest.Manifest
	handlers map[string]Handler
}

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
	Status  int
	Body    map[string]interface{}
	Headers map[string][]string
}

type Handler func(*Request, *Response)

func (s *Service) Get(path string, handler Handler) {
	s.handle("GET", path, handler)
}

func (s *Service) Put(path string, handler Handler) {
	s.handle("PUT", path, handler)
}

func (s *Service) Post(path string, handler Handler) {
	s.handle("POST", path, handler)
}

func (s *Service) Delete(path string, handler Handler) {
	s.handle("DELETE", path, handler)
}

func (s *Service) handle(method, path string, handler Handler) {
	handlerId := uuid.NewV4().String()
	s.manifest.Routes = append(s.manifest.Routes, &manifest.Route{
		HttpMethod: manifest.HttpMethodFromString(method),
		PathSpec:   path,
		Handler:    handlerId,
	})
	s.handlers[handlerId] = handler
}
