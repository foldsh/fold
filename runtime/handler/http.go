package handler

import (
	"context"
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

func (h *HTTPHandler) Serve() error {
	if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *HTTPHandler) Shutdown(ctx context.Context, done chan struct{}) {
	if err := h.server.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		h.logger.Errorf("HTTP server Shutdown: %v", err)
	}
	close(done)
}
