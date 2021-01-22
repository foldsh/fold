package main

import (
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
)

func main() {
	logger, err := logging.NewLogger(logging.Debug, false)
	if err != nil {
		panic("failed to start logger")
	}
	runtime.HTTP(logger, os.Args[1], os.Args[2:]...)
}
