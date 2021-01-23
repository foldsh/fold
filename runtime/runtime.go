package runtime

import (
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/router"
	"github.com/foldsh/fold/runtime/supervisor"
)

var rt *runtime

func HTTP(logger logging.Logger, command string, args ...string) {
	rt = initRuntime(logger, command, args...)
	rt.handler = handler.NewHTTP(logger, rt.router, ":8080")
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
	r.handler.Serve()
}
