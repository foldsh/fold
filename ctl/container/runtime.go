package container

import (
	"context"
	"io"

	"github.com/docker/docker/client"
	"github.com/foldsh/fold/logging"
)

func NewRuntime(
	ctx context.Context, logger logging.Logger, out io.Writer,
) (*ContainerRuntime, error) {
	client, err := newDockerClient(logger)
	if err != nil {
		return nil, err
	}
	return &ContainerRuntime{
		cli:    client,
		ctx:    ctx,
		logger: logger,
		out:    out,
	}, nil
}

type ContainerRuntime struct {
	cli    DockerClient
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
}

func (cr *ContainerRuntime) Context() context.Context {
	return cr.ctx
}

func newDockerClient(logger logging.Logger) (DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Debugf("failed to initialised docker client")
		return nil, FailedToConnectToDockerEngineError
	}
	return client, nil
}
