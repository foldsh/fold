// Package runtime manages all of the components required to run a users application. It is
// implemented as a state machine. This provides an easy way to manage the relationship between
// the state of the underlying process and the behaviour of the runtime.
package runtime

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
	"github.com/foldsh/fold/runtime/fsm"
	"github.com/foldsh/fold/runtime/router"
	"github.com/foldsh/fold/runtime/supervisor"
	"github.com/foldsh/fold/runtime/transport"
)

type Supervisor interface {
	Start(env map[string]string) error
	Restart(env map[string]string) error
	Stop() error
	Kill() error
	Wait() error
	Signal(sig os.Signal) error
}

type Client interface {
	Start(string) error
	Stop() error
	Restart(string) error
	GetManifest(context.Context) (*manifest.Manifest, error)
	DoRequest(context.Context, *transport.Request) (*transport.Response, error)
}

type Router interface {
	http.Handler
	Configure(*manifest.Manifest)
}

type SocketFactory func() string

type RouterFactory func(logger logging.Logger, doer router.RequestDoer) Router

type Runtime struct {
	logger logging.Logger
	fsm    *fsm.FSM
	cmd    string
	args   []string
	done   chan struct{}

	// These properties have the same lifetime as the runtime
	env           map[string]string
	supervisor    Supervisor
	client        Client
	socketFactory SocketFactory
	routerFactory RouterFactory
	defaultRouter Router
	onProcessEnd  func()

	// These are set dynamically with restarts etc
	socketAddress string
	router        Router
}

var (
	UP     fsm.State = "UP"
	DOWN   fsm.State = "DOWN"
	EXITED fsm.State = "EXITED"

	START       fsm.Event = "START"
	STOP        fsm.Event = "STOP"
	EXIT        fsm.Event = "EXIT"
	CRASH       fsm.Event = "CRASH"
	FILE_CHANGE fsm.Event = "FILE_CHANGE"
)

func NewRuntime(
	logger logging.Logger,
	cmd string,
	args []string,
	done chan struct{},
	options ...Option,
) *Runtime {
	newRuntime := &Runtime{
		logger: logger,
		cmd:    cmd,
		args:   args,
		done:   done,
	}

	// First up we configure the default FSM. Other options can change it later on.
	configureFSM(newRuntime)

	// The default options are handled the same way as user defined options. Options are applied
	// in order so the defaults just get overriden by the user defined ones.
	defaultOptions := []Option{
		WithSupervisor(
			supervisor.NewSupervisor(
				newRuntime.logger,
				newRuntime.cmd,
				newRuntime.args,
				os.Stdout,
				os.Stdout,
			),
		),
		WithClient(transport.NewIngress(newRuntime.logger)),
		WithSocketFactory(newAddr),
		WithRouterFactory(func(l logging.Logger, d router.RequestDoer) Router {
			return router.NewRouter(l, d)
		}),
		WithDefaultRouter(router.NewCatchAllRouter(newRuntime.logger, &defaultRequestDoer{})),
		// For now, regardless of the reason for termination, we handle process termination using
		// a CRASH event. This is because we currently only support long lived processes like
		// servers which are terminated from the outside. When we support batch jobs this will
		// change.
		OnProcessEnd(func() { newRuntime.Emit(CRASH) }),
	}

	// Then we go through the options specified by the caller and apply all of them
	for _, option := range append(defaultOptions, options...) {
		option(newRuntime)
	}

	return newRuntime
}

// For now, if we encounter any errors starting or stopping the process we exit.
// This does not occur when the users process crashes or fails to start. That shows up as a
// successful process start and a crash.
func configureFSM(r *Runtime) {
	f := fsm.NewFSM(
		r.logger,
		DOWN,
		fsm.Transitions{
			{START, DOWN, UP, []fsm.Callback{
				func() {
					if err := r.startClientAndSupervisor(); err != nil {
						r.exit()
					}
				},
			}},
			{START, UP, UP, []fsm.Callback{
				func() {
					if err := r.restartClientAndSupervisor(); err != nil {
						r.exit()
					}
				},
			}},
			{STOP, UP, EXITED, []fsm.Callback{
				func() {
					if err := r.stopClientAndSupervisor(); err != nil {
						r.exit()
					}
				}},
			},
			{STOP, DOWN, EXITED, nil},
			{CRASH, UP, EXITED, nil},
			{EXIT, UP, EXITED, nil},
			{EXIT, DOWN, EXITED, nil},
		},
	)
	// Whenever we transition back to DOWN or EXIT, we want to switch back to the default router.
	// Transitioning to EXITED will result in a shutdown pretty snappily but we'll set the
	// default router up again so there is a semblance of graceful handling.
	f.OnTransitionTo(DOWN, func() { r.router = r.defaultRouter })
	f.OnTransitionTo(EXITED, func() {
		r.router = r.defaultRouter
		close(r.done)
	})

	// Finally we set the FSM on the runtime
	r.fsm = f
}

func (r *Runtime) State() fsm.State {
	return r.fsm.State()
}

func (r *Runtime) Router() Router {
	return r.router
}

func (r *Runtime) Start() {
	r.Emit(START)
}

func (r *Runtime) Stop() {
	r.Emit(STOP)
}

func (r *Runtime) Signal(signal os.Signal) {
	r.supervisor.Signal(signal)
}

func (r *Runtime) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Runtime) Emit(event fsm.Event) {
	r.fsm.Emit(event)
}

func (r *Runtime) createAndConfigureRouter() error {
	r.logger.Debugf("Setting up new router")
	r.router = r.routerFactory(r.logger, r.client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	manifest, err := r.client.GetManifest(ctx)
	if err != nil {
		r.logger.Debugf("failed to fetch manifest")
		return err
	}
	r.router.Configure(manifest)
	return nil
}

func (r *Runtime) startClientAndSupervisor() error {
	r.logger.Debugf("Starting the client and supervisor")
	r.socketAddress = r.socketFactory()
	env := map[string]string{"FOLD_SOCK_ADDR": r.socketAddress}
	if err := r.supervisor.Start(env); err != nil {
		return err
	}
	// Now that we've started the process, we want to set up a goroutine that waits for the process
	// to terminate. When it does, we'll identify whether it was a crash or not and then emit
	// the appropriate event.
	go func() {
		err := r.supervisor.Wait()
		// If the process was stopped by a signal then it was intentional and we want to exit,
		// regardless of the configured process end behaviour.
		if errors.Is(err, supervisor.TerminatedBySignal) {
			r.stopClientAndSupervisor()
			r.exit()
			return
		}
		r.onProcessEnd()
	}()
	if err := r.client.Start(r.socketAddress); err != nil {
		return err
	}
	if err := r.createAndConfigureRouter(); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) stopClientAndSupervisor() error {
	r.logger.Debugf("Stopping the client and supervisor")
	if err := r.client.Stop(); err != nil {
		return err
	}
	if err := r.supervisor.Stop(); err != nil {
		return err
	}
	// It's important to note that we don't wait here. Waiting is handled by the goroutine that
	// is started when the process itself is started. Waiting twice can lead to confusion about
	// where the notification comes in and how to handle it, so it's better just to do it in one
	// place.
	return nil
}

func (r *Runtime) restartClientAndSupervisor() error {
	if err := r.stopClientAndSupervisor(); err != nil {
		return err
	}
	if err := r.startClientAndSupervisor(); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) exit() {
	r.Emit(EXIT)
}

type defaultRequestDoer struct{}

func (d *defaultRequestDoer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte(`{"title":"service is down"}`))
}
