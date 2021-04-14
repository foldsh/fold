package container_test

import (
	"context"
	"os"

	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/mocks"
	"github.com/foldsh/fold/logging"
)

func setup() (*container.ContainerRuntime, *mocks.DockerClient, *mocks.FileSystem) {
	dc := &mocks.DockerClient{}
	fs := &mocks.FileSystem{}
	rt, _ := container.NewRuntime(
		context.Background(),
		logging.NewTestLogger(),
		os.Stdout,
		fs,
		dc,
	)
	return rt, dc, fs
}
