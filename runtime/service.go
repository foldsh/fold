/**
- Separate out the handler from the router. The router is internal to the runtime but the handler
is external.
- The runtime should just expose one method for the handler to invoke with an HTTP request.
- This lets us issue an event for each request so that we hook other things into it.

- Requests in particular should not be serialised. There should be a direct method call from the
goroutine which is running the handler straight into the gRPC stack. The gRPC client is thread safe
and can be used to handle multiple requests in parallel.

- An alternative implementation is to buffer all of the requests in a queue and then have multiple
workers 'reparallelise' them before hitting the gRPC client. This sounds redundant at first but would
give the runtime fine control over exactly how many requests made it through to the users code.
This could be used to provide adaptive backpressure based on queue statistics, providing a trigger
for scaling up/down. Don't bother with this for now as the implementation will all be lambda based
at first anyway, but include it in some commentary or write it up in notion.

- The runtime behaviour should be configurable by options. These should result in certain functionality
being hooked on certain events.

- Registering hot reload or restart-on-crash handlers should result in requests being given default response
while the process is down. We therefore need a concept of runtime state in addition to hooks for handlers.
Access to that state will need to be serialised but that is ok in a dev environment. That penalty should
only be incurred if the restart type handlers are registered

- We won't need many event types at first but it should give a pretty flexible design. It also keeps
all of the different bits quite nicely separated as you can just test that the handlers are invoked
for a given event.

- This also calls for an integration test of sorts. Look into setting that up. Should it be part of
this runtime package? With a 'slow' flag? Or part of a top level 'tests' package? Would it be
possible to use 'go run' as the command? Then we could actually use the go SDK quite easily.
nomad appears to be quite a good example. They just structure it as a regular go package at the top
level. They have built a little framework for their use case but ordinary tests would be fine for
me for now.

QUESTIONS
 - We want to make the client and supervisor parameters to the runtime so they can be mocked. However, we also actually need a way to configure those objects without the constructor.
 This is because on a restart we need to create a new client and a new process.
   - one option is to pass factory methods as dependencies
   - another option is to make the objects reconfigurable
   - For now, lets pass a client factory but make the supervisor configure the env
     on start/restart

*/
package runtime

import (
	"context"
	"net/http"
	"os"

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
	CRASH       fsm.Event = "CRASH"
	FILE_CHANGE fsm.Event = "FILE_CHANGE"
)

// This can probably be modelled pretty well as a state machine. Probably
// don't want to bring in a library for it though.
// In the 'DOWN' state, only a 'START' event does anything, it configures a
// process/client etc. A default router handles requests in this state and it
// just gives 500 errors.
//
// Failing to START leaves us in DOWN.
// A successful START leaves us in UP. In this state requests are handled
// by the real router that has been configured from the running process.
//
// The two events that happen in the UP state are 'CRASH' and 'START'.
// In the UP state, a START event behaves like a restart. It will first shut
// the process down, and then cause a restart. Note that this handler is only
// registered if the relevant option is passed. Therefore, in prod mode, there
// is no START handler for the UP state.
//
// The CRASH event is handled differently depending on how it was configured.
// by default, the CRASH event will simply result in flushing the logs for
// the process and then exiting.
// It can be configured, however to simply put the runtime back into the DOWN
// state. This will set it up to be restarted by a FILE_CHANGE
//
// The FILE_CHANGE event is only emitted if a watcher has been registered
// through the appropriate option. The handler for it simply emits a
// a START event. This will have the effect of restarting if that is
// how the runtime is configured, or just starting if it is currently DOWN.
func NewRuntime(
	logger logging.Logger,
	Cmd string,
	Args []string,
	options ...Option,
) *Runtime {
	newRuntime := &Runtime{
		logger: logger,
		cmd:    Cmd,
		args:   Args,
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

func configureFSM(r *Runtime) {
	f := fsm.NewFSM(
		r.logger,
		DOWN,
		fsm.Transitions{
			{START, DOWN, UP, []fsm.Callback{
				func() {
					// TODO how do we handle the errors? Need to add an error
					// into the Callback type?
					if err := r.startClientAndSupervisor(); err != nil {
						return
					}
				},
			}},
			{START, UP, UP, []fsm.Callback{
				func() {
					if err := r.restartClientAndSupervisor(); err != nil {
						// TODO how to handle this? This should probably result in a transition
						// back to the DOWN state, or EXIT, depending on how CRASH is handled.
						return
					}
				},
			}},
			{STOP, UP, EXITED, []fsm.Callback{
				func() {
					if err := r.stopClientAndSupervisor(); err != nil {
						// TODO how to handle this? This should probably result in a transition
						// back to the DOWN state, or EXIT, depending on how CRASH is handled.
						return
					}
				}},
			},
			{STOP, DOWN, EXITED, nil},
			{CRASH, UP, EXITED, nil},
		},
	)
	// Whenever we transition back to DOWN or EXIT, we want to switch back to the default router.
	// Transitioning to EXITED will result in a shutdown pretty snappily but we'll set the
	// default router up again so there is a semblance of graceful handling.
	f.OnTransitionTo(DOWN, func() { r.router = r.defaultRouter })
	f.OnTransitionTo(EXITED, func() { r.router = r.defaultRouter })

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
	// TODO context should have a timeout
	manifest, err := r.client.GetManifest(context.Background())
	if err != nil {
		r.logger.Fatalf("failed to fetch manifest")
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
		r.supervisor.Wait()
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

// TODO it will probably be convenient to make the Stop/Kill methods etc idempotent, or return
// a useful error. I.e. callng stop on a stopped supervisor/client should be safe.
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

type defaultRequestDoer struct{}

func (d *defaultRequestDoer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte(`{"title":"service is down"}`))
}
