package main

import (
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/service"
)

func main() {
	logger, err := logging.NewLogger(logging.Debug, false)
	if err != nil {
		panic("failed to start logger")
	}
	cmd := service.Command{os.Args[1], os.Args[2:]}
	logger.Debugf("starting service with command %+v", cmd)
	service, err := service.NewService(logger, cmd)
	if err != nil {
		panic(err)
	}
	handler := handler.NewHTTP(logger, service)
	handler.Start()
}
