package runtime

import (
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/fsm"
	"github.com/foldsh/fold/runtime/watcher"
)

type Option func(*Runtime)

type HandlerT uint8

func WithEnv(env map[string]string) Option {
	return func(r *Runtime) {}
}

func WithSupervisor(supervisor Supervisor) Option {
	return func(r *Runtime) {
		r.supervisor = supervisor
	}
}

func WithClient(client Client) Option {
	return func(r *Runtime) {
		r.client = client
	}
}

func WithSocketFactory(socketFactory SocketFactory) Option {
	return func(r *Runtime) {
		r.socketFactory = socketFactory
	}
}

func WithRouterFactory(routerFactory RouterFactory) Option {
	return func(r *Runtime) {
		r.routerFactory = routerFactory
	}
}

func WithDefaultRouter(router Router) Option {
	return func(r *Runtime) {
		r.defaultRouter = router
		r.router = r.defaultRouter
	}
}

func OnProcessEnd(handler func()) Option {
	return func(r *Runtime) {
		r.onProcessEnd = handler
	}
}

func WatchDir(frequency time.Duration, dir string) Option {
	return func(r *Runtime) {
		debouncer := watcher.NewDebouncer(frequency, func() {
			r.Emit(FILE_CHANGE)
		})
		watcher, err := watcher.NewWatcher(r.logger, dir, debouncer.OnChange)
		if err != nil {
			r.logger.Fatalf("Failed to setup hot reloading")
		}
		if err := watcher.Watch(); err != nil {
			r.logger.Fatalf("Failed to setup hot reloading")
		}
		r.fsm.AddTransition(fsm.Transition{FILE_CHANGE, UP, UP, []fsm.Callback{
			func() {
				if err := r.restartClientAndSupervisor(); err != nil {
					return
				}
			},
		}})
		r.fsm.AddTransition(fsm.Transition{FILE_CHANGE, DOWN, UP, []fsm.Callback{
			func() {
				if err := r.startClientAndSupervisor(); err != nil {
					return
				}
			},
		}})
		r.fsm.OnTransitionTo(EXITED, func() {
			watcher.Close()
		})
	}
}

type CrashPolicyT uint8

const (
	KILL CrashPolicyT = iota + 1
	KEEP_ALIVE
)

func CrashPolicy(crashPolicy CrashPolicyT) Option {
	return func(r *Runtime) {
		switch crashPolicy {
		case KILL:
			// This is the default setting so we don't need to do anything.
			return
		case KEEP_ALIVE:
			r.fsm.AddTransition(fsm.Transition{CRASH, UP, DOWN, nil})
		}
	}
}

func LogLevel(level logging.LogLevel) Option {
	return func(r *Runtime) {}
}
