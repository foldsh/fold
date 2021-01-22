package handler

import (
	"net/http"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

type ConfigurableHTTPServer interface {
	http.Handler
	Configure(*manifest.Manifest)
}

func NewHTTP(logger logging.Logger, server ConfigurableHTTPServer, port string) Handler {
	return &httpHandler{logger, server, port}
}

type httpHandler struct {
	logger logging.Logger
	server ConfigurableHTTPServer
	port   string
}

func (hh *httpHandler) Serve() {
	http.ListenAndServe(hh.port, hh.server)
}
