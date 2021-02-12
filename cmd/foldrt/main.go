package main

import (
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
)

func main() {
	var (
		logger logging.Logger
		err    error
	)
	stage := os.Getenv("FOLD_STAGE")

	if stage == "TEST" {
		logger, err = logging.NewLogger(logging.Debug, true)
	} else if stage == "PROD" {
		logger, err = logging.NewLogger(logging.Info, true)
	} else {
		logger, err = logging.NewLogger(logging.Debug, false)
	}
	logger.Debug("starting fold runtime for stage ", stage)

	if err != nil {
		panic("failed to start logger")
	}

	env := os.Getenv("FOLD_ENV")
	runtime.Run(logger, env, stage, os.Args[1], os.Args[2:]...)
}
