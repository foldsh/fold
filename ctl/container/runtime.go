package container

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/foldsh/fold/logging"
)

func NewRuntime(
	ctx context.Context,
	logger logging.Logger,
	out io.Writer,
	fs FileSystem,
	client DockerClient,
) (*ContainerRuntime, error) {
	return &ContainerRuntime{
		cli:    client,
		ctx:    ctx,
		logger: logger,
		out:    out,
		fs:     osFileSystem{},
	}, nil
}

type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
}

type ContainerRuntime struct {
	cli    DockerClient
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
	fs     FileSystem
}

func (cr *ContainerRuntime) Context() context.Context {
	return cr.ctx
}

func NewDockerClient(logger logging.Logger) (DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Debugf("failed to initialised docker client")
		return nil, FailedToConnectToDockerEngineError
	}
	return client, nil
}

type osFileSystem struct{}

func (fs osFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
