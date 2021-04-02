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

	env           map[string]string
	supervisor    Supervisor
	client        Client
	socketFactory SocketFactory
	routerFactory RouterFactory

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
	r := &Runtime{
		logger: logger,
		cmd:    Cmd,
		args:   Args,
	}

	// First up we go through the options specified by the caller and apply all of them
	for _, option := range options {
		option(r)
	}

	// Then we go through and set the defaults where they are needed.
	// TODO move this to before the options and set up the default FSM in there.
	setDefaultOptions(r)

	// Next we set up the FSM
	f := fsm.NewFSM(
		DOWN,
		fsm.Transitions{
			{START, DOWN, UP, []fsm.Callback{
				func() {
					// TODO how do we handle the errors? Need to add an error
					// into the Callback type?
					r.socketAddress = r.socketFactory()
					env := map[string]string{"FOLD_SOCK_ADDR": r.socketAddress}
					if err := r.supervisor.Start(env); err != nil {
						return
					}
					if err := r.client.Start(r.socketAddress); err != nil {
						return
					}
					return
				}},
			},
			{STOP, DOWN, EXITED, nil},
			{STOP, UP, EXITED, nil},
			// TODO when the appropriate option is set we need to choose one of these two:
			{CRASH, UP, EXITED, nil},
			{CRASH, UP, DOWN, nil},
			// TODO when the option is set, this needs to be registered
			{FILE_CHANGE, UP, DOWN, []fsm.Callback{
				func() {
					// TODO this needs to reference the fsm so that it can
					// emit the start event when it's done.
					// f.Emit(START)
					// We need to provide a way to build up an FSM in multiple
					// method calls basically.
				},
			}},
		},
	)
	r.fsm = f
	return r
}

func setDefaultOptions(r *Runtime) {
	if r.supervisor == nil {
		r.supervisor = supervisor.NewSupervisor(r.logger, r.cmd, r.args, os.Stdout, os.Stdout)
	}
	if r.client == nil {
		r.client = transport.NewIngress(r.logger)
	}
	if r.socketFactory == nil {
		r.socketFactory = newAddr
	}
	if r.routerFactory == nil {
		r.routerFactory = func(l logging.Logger, d router.RequestDoer) Router {
			return router.NewRouter(l, d)
		}
	}
}

func (r *Runtime) State() fsm.State {
	return r.fsm.State()
}

func (r *Runtime) Router() Router {
	return r.router
}

func (r *Runtime) Start() {
	r.Trigger(START)
}

func (r *Runtime) Stop() {
	r.Trigger(STOP)
}

func (r *Runtime) ServeHTTP(http.ResponseWriter, *http.Request) {

}

func (r *Runtime) Trigger(event fsm.Event) {
	r.fsm.Emit(event)
}

func (r *Runtime) configure() error {
	r.logger.Debugf("Fetching manifest")
	// TODO context
	manifest, err := r.client.GetManifest(context.Background())
	if err != nil {
		r.logger.Fatalf("failed to fetch manifest")
	}
	r.router.Configure(manifest)
	return nil
}
