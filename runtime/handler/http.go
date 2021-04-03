package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/foldsh/fold/logging"
)

func NewHTTP(
	logger logging.Logger,
	handler http.Handler,
	addr string,
) *HTTPHandler {
	return &HTTPHandler{logger: logger, server: &http.Server{Addr: addr, Handler: handler}}
}

type HTTPHandler struct {
	logger logging.Logger
	server *http.Server
}

func (h *HTTPHandler) Serve() {
	if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func (h *HTTPHandler) Shutdown() {
	if err := h.server.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		h.logger.Errorf("HTTP server Shutdown: %v", err)
	}
}
