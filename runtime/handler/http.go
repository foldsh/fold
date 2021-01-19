package handler

import (
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/service"
)

type httpHandler struct {
	logger  logging.Logger
	service service.Service
}

func NewHTTP(logger logging.Logger, service service.Service) Handler {
	return &httpHandler{logger, service}
}

func (hh *httpHandler) Start() {
	hh.service.Start()
}
