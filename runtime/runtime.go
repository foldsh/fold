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
	"github.com/foldsh/fold/runtime/subprocess"
	"github.com/foldsh/fold/runtime/watcher"
)

var (
	loggr logging.Logger
	env   string
	stage string
	suprv *supervisor.Supervisor
	mnfst *manifest.Manifest
	routr router.Router
	hndlr Handler
)

type Handler interface {
	Serve()
}

func Run(logger logging.Logger, env string, stage string, command string, args ...string) {
	// TODO this is optimised for local development... it's a bit convoluted but
	// it basically tries really hard to always be available. This is so that
	// the app doesn't constantly crash for local development when errors are common.
	// Ideally an error should just result in a helpful message and then as soon as
	// the fix is written to disk we try to load the server again.
	// For deployment purposes though, whether to a test or prodution environment,
	// it would be much better to have the whole thing just fail as fast as possible.
	// I'll split this up into two mains I think to accomplish this.
	start := time.Now()
	loggr = logger
	env = env
	stage = stage
	suprv = supervisor.NewSupervisor(loggr, command, args...)
	routr = router.NewRouter(loggr, suprv)

	// Start watching for file change immediately so that we can recover from
	// someone overwriting the command by binding an empty directory.
	if stage == "LOCAL" {
		setupHotReload()
	}

	registerSignalHandlers()

	// Setup handler
	switch env {
	case "LAMBDA":
		hndlr = handler.NewLambda(loggr, routr)
	default:
		hndlr = handler.NewHTTP(loggr, routr, ":6123")
	}
	go func() {
		// We start the handler in a goroutine like this to ensure that
		// it always comes up and is responsive.
		hndlr.Serve()
	}()

	// Start supervisor
	loggr.Debugf("starting supervisor")
	err := suprv.Start()
	if err != nil {
		loggr.Fatalf("supervisor failed to start process")
	}

	fetchManifestAndConfigureRouter()

	elapsed := time.Since(start)
	loggr.Infof("ready to accept requests, startup took %s", elapsed)
	done := make(chan struct{})
	<-done
}

func registerSignalHandlers() {
	loggr.Debugf("registering signal handlers")
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		err := suprv.Signal(s)
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()
}

func fetchManifestAndConfigureRouter() {
	loggr.Debugf("fetching manifest")
	mnfst, err := suprv.GetManifest()
	if err != nil {
		loggr.Fatalf("failed to fetch manifest")
	}
	loggr.Debugf("router is %+v", routr)
	routr.Configure(mnfst)
}

func setupHotReload() {
	watchdir := os.Getenv("FOLD_WATCH_DIR")
	loggr.Debugf("watching for changes in %s", watchdir)
	if watchdir == "" {
		return
	}
	reloadFn := func() {
		if err := suprv.Restart(); err != nil {
			loggr.Fatalf("failed to restart")
		}
		fetchManifestAndConfigureRouter()
	}
	debouncer := watcher.NewDebouncer(100*time.Millisecond, reloadFn)
	watcher, err := watcher.NewWatcher(loggr, watchdir, debouncer.OnChange)
	if err != nil {
		loggr.Fatalf("failed to setup hot reloading")
	}
	if err := watcher.Watch(); err != nil {
		loggr.Fatalf("failed to setup hot reloading")
	}
}
