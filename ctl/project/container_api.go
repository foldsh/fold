package project

import (
	"context"

	"github.com/foldsh/fold/ctl/container"
)

type ContainerAPI interface {
	Context() context.Context
	NewNetwork(name string) *container.Network
	NetworkExists(net *container.Network) (bool, error)
	CreateNetwork(net *container.Network) error
	RemoveNetwork(net *container.Network) error
	GetImage(name string) (*container.Image, error)
	PullImage(name string) (*container.Image, error)
	BuildImage(img *container.Image) error
	NewContainer(
		name string,
		image container.Image,
		mounts ...container.Mount,
	) *container.Container
	GetContainer(name string) (*container.Container, error)
	RunContainer(net *container.Network, con *container.Container) error
	StopContainer(con *container.Container) error
	ContainerLogs(con *container.Container) (*container.LogStream, error)
}
