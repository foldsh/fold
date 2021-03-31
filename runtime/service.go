package runtime

import (
	"context"
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/types"
)

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

type RuntimeOpts struct {
	Cmd        string
	Args       []string
	Env        map[string]string
	Supervisor Supervisor
	Client     Client
}

func NewRuntime(logger logging.Logger, opts RuntimeOpts) *Runtime {
	return &Runtime{
		logger:     logger,
		cmd:        opts.Cmd,
		args:       opts.Args,
		supervisor: opts.Supervisor,
		client:     opts.Client,
	}
}
