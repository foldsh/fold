package handler

import (
	"net/http"
	"os"

	"github.com/foldsh/fold/logging"
)

func NewHTTP(logger logging.Logger, server http.Handler, port string) *HTTPHandler {
	return &HTTPHandler{logger, server, port}
}

type HTTPHandler struct {
	logger logging.Logger
	server http.Handler
	port   string
}

func (hh *HTTPHandler) Serve() {
	hh.logger.Debugf("Listening on port %s", hh.port)
	if err := http.ListenAndServe(hh.port, hh.server); err != nil {
		hh.logger.Errorf(err.Error())
		os.Exit(1)
	}
}
