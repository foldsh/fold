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
*/
package runtime

import (
	"context"
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/types"
)

type Handler interface {
	Serve()
}

type Supervisor interface {
	Start() error
	Restart() error
	Stop() error
	Kill() error
	Wait() error
	Signal(sig os.Signal) error
}

type Client interface {
	Start() error
	GetManifest(ctx context.Context) error
	DoRequest(ctx context.Context, req *types.Request) (*types.Response, error)
}

type Runtime struct {
	logger     logging.Logger
	cmd        string
	args       []string
	env        map[string]string
	supervisor Supervisor
	client     Client

	socketAddress string
}

type EventT int

const (
	Start EventT = iota + 1
	Stop
	Crash
	FileChange
)

type EventHandler func()

type RuntimeOpts struct {
	Cmd        string
	Args       []string
	Env        map[string]string
	Supervisor Supervisor
	Client     Client

	handlers map[EventT][]EventHandler
}

func NewRuntime(logger logging.Logger, opts RuntimeOpts) *Runtime {
	return &Runtime{
		logger:     logger,
		cmd:        opts.Cmd,
		args:       opts.Args,
		supervisor: opts.Supervisor,
		client:     opts.Client,
		handlers:   map[EventT]EventHandler{},
	}
}

func (r *Runtime) Start() error {

}

func (r *Runtime) Configure() error {
	r.logger.Debugf("Fetching manifest")
	manifest, err := r.client.GetManifest()
	if err != nil {
		r.logger.Fatalf("failed to fetch manifest")
	}
	loggr.Debugf("router is %+v", routr)
	routr.Configure(mnfst)
}

func (r *Runtime) Stop() error {

}

func (r *Runtime) DoRequest(http.ResponseWriter, *http.Request) {

}

func (r *Runtime) subscribe(event EventT, handler EventHandler) {
	handlers = r.handlers[event]
	r.handlers[event] = append(handlers, handler)
}

func (r *Runtime) publish(event EventT) {
	for _, handler := range r.handlers[event] {
		handler()
	}
}
