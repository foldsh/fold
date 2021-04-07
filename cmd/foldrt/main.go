package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	handlerImpl "github.com/foldsh/fold/runtime/handler"
)

type Handler interface {
	Serve() error
	Shutdown(context.Context, chan struct{})
}

func main() {
	var (
		logger  logging.Logger
		err     error
		options []runtime.Option
		handler Handler
	)

	// TEST, PROD
	stage := os.Getenv("FOLD_STAGE")
	// LAMBDA, HTTP
	env := os.Getenv("FOLD_ENV")
	watchDir := os.Getenv("FOLD_WATCH_DIR")

	switch stage {
	case "TEST":
		logger, err = logging.NewLogger(logging.Debug, true)
	case "PROD":
		logger, err = logging.NewLogger(logging.Info, true)
	default:
		// Local development mode
		logger, err = logging.NewLogger(logging.Debug, false)
		options = append(options, runtime.CrashPolicy(runtime.KEEP_ALIVE))
		if watchDir != "" {
			options = append(options, runtime.WatchDir(100*time.Millisecond, watchDir))
		}
	}

	if err != nil {
		panic("failed to start logger")
	}

	logger.Debug("Starting fold runtime for stage: ", stage)

	runtimeStopped := make(chan struct{})
	rt := runtime.NewRuntime(logger, os.Args[1], os.Args[2:], runtimeStopped, options...)
	rt.Start()

	switch env {
	case "LAMBDA":
		handler = handlerImpl.NewLambda(logger, rt)
	default:
		handler = handlerImpl.NewHTTP(logger, rt, ":6123")
	}

	handlerShutdown := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		s := <-signals

		// Ok we got a signal to kill the application. First we shutdown the server gracefully:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		handler.Shutdown(ctx, handlerShutdown)
		// Now we can signal the runtime.
		rt.Signal(s)
	}()
	if err := handler.Serve(); err != nil {
		// Error starting or closing listener:
		logger.Fatalf("Server stopped unexpectedly: %v", err)
	}
	logger.Debugf("waiting for handler shutdown")
	<-handlerShutdown
	logger.Debugf("waiting for runtime to stop")
	<-runtimeStopped
}
