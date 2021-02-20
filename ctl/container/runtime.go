package container

import (
	"context"
	"io"
	"os"

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
		fs:     osFileSystem{},
	}, nil
}

type ContainerRuntime struct {
	cli    DockerClient
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
	fs     fileSystem
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

type fileSystem interface {
	mkdirAll(path string, perm os.FileMode) error
}

type osFileSystem struct{}

func (fs osFileSystem) mkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
