package handler

import (
	"net/http"

	"github.com/foldsh/fold/logging"
)

func NewHTTP(logger logging.Logger, server http.Handler, port string) Handler {
	return &httpHandler{logger, server, port}
}

type httpHandler struct {
	logger logging.Logger
	server http.Handler
	port   string
}

func (hh *httpHandler) Serve() {
	http.ListenAndServe(hh.port, hh.server)
}
