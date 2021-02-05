package container

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/foldsh/fold/logging"
)

var (
	DockerEngineError                  = errors.New("failed to connect to docker engine")
	FailedToConnectToDockerEngineError = errors.New("failed to connect to docker engine")
	FailedToPrepareBuildArchive        = errors.New("failed to prepare the build archive")
	FailedToBuildImage                 = errors.New("failed to build the image")
	FailedToPullImage                  = errors.New("failed to pull the image")
	FailedToCreateContainer            = errors.New("failed to create the container")
	FailedToStartContainer             = errors.New("failed to start the container")
	FailedToStopContainer              = errors.New("failed to stop the container")
	FailedToRemoveContainer            = errors.New("failed to remove the container")

	foldPrefix = "fold."
)

type ContainerRuntime interface {
	NewContainer(name string, image string, volumes ...Volume) *Container
	NewNetwork(name string) *Network
	AllContainers() ([]*Container, error)
	GetContainer(name string) (*Container, error)
}

func NewRuntime(
	ctx context.Context, logger logging.Logger, out io.Writer, client DockerClient,
) ContainerRuntime {
	return &containerRuntime{
		cli:    client,
		ctx:    ctx,
		logger: logger,
		out:    out,
	}
}

type containerRuntime struct {
	cli    DockerClient
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
}

func (cr *containerRuntime) NewContainer(
	name string, image string, volumes ...Volume,
) *Container {
	return &Container{
		Name:    fmt.Sprintf("%s%s", foldPrefix, name),
		Image:   image,
		Volumes: volumes,
		rt:      cr,
	}
}

func (cr *containerRuntime) NewNetwork(name string) *Network {
	return &Network{Name: name, rt: cr}
}

func (cr *containerRuntime) AllContainers() ([]*Container, error) {
	return cr.listContainers()
}

func (cr *containerRuntime) GetContainer(name string) (*Container, error) {
	name = fmt.Sprintf("%s%s", foldPrefix, name)
	containers, err := cr.AllContainers()
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, nil
}

// A very simple struct for representing the functionality we actually
// need to expose from here.
type Container struct {
	ID      string
	Name    string
	Image   string
	Volumes []Volume

	rt *containerRuntime
}

func (c *Container) Start() error {
	return c.rt.runContainer(c)
}

func (c *Container) Stop() error {
	return c.rt.stopContainer(c)
}

func (c *Container) Remove() error {
	return c.rt.removeContainer(c)
}

func (c *Container) JoinNetwork(net *Network) error {
	return c.rt.addToNetwork(net, c)
}

func (c *Container) LeaveNetwork(net *Network) error {
	return c.rt.removeFromNetwork(net, c)
}

type Volume struct {
	Src string
	Dst string
}
