package runtime

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/router"
	"github.com/foldsh/fold/runtime/supervisor"
)

var rt *runtime

func HTTP(logger logging.Logger, command string, args ...string) {
	start := time.Now()
	rt = initRuntime(logger, command, args...)
	rt.handler = handler.NewHTTP(logger, rt.router, ":8080")
	elapsed := time.Since(start)
	logger.Infof("ready to accept requests, startup took %s", elapsed)
	rt.Start()
}

func Lambda(logger logging.Logger, command string, args ...string) {
	rt = initRuntime(logger, command, args...)
	rt.handler = handler.NewLambda(logger, rt.router)
	rt.Start()
}

type runtime struct {
	logger   logging.Logger
	superv   supervisor.Supervisor
	manifest *manifest.Manifest
	router   router.Router
	handler  handler.Handler
}

func initRuntime(logger logging.Logger, command string, args ...string) *runtime {
	superv := supervisor.NewSupervisor(logger)
	err := superv.Exec(command, args...)
	if err != nil {
		logger.Fatalf("supervisor failed to start subprocess")
	}
	manifest, err := superv.GetManifest()
	if err != nil {
		logger.Fatalf("failed to fetch manifest")
	}
	router := router.NewRouter(logger, superv)
	router.Configure(manifest)
	return &runtime{logger: logger, superv: superv, manifest: manifest, router: router}
}

func (r *runtime) Start() {
	r.registerSignalHandlers()
	r.handler.Serve()
}

func (r *runtime) registerSignalHandlers() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		err := r.superv.Signal(s)
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()
}
