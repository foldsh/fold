package runtime

import (
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/router"
	"github.com/foldsh/fold/runtime/supervisor"
)

func HTTP(logger logging.Logger, command string, args ...string) {
	supervisor := supervisor.NewSupervisor(logger)
	err := supervisor.Exec(command, args...)
	if err != nil {
		logger.Fatalf("supervisor failed to start subprocess")
	}
	manifest, err := supervisor.GetManifest()
	if err != nil {
		logger.Fatalf("failed to fetch manifest")
	}
	router := router.NewRouter(logger, supervisor)
	router.Configure(manifest)
	handler := handler.NewHTTP(logger, router, ":8080")
	handler.Serve()
}
