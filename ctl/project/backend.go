package project

import (
	"context"

	"github.com/foldsh/fold/ctl/container"
)

type Backend interface {
	Context() context.Context
	NewNetwork(name string) *container.Network
	CreateNetworkIfNotExists(net *container.Network) error
	RemoveNetworkIfExists(net *container.Network) error
	NewContainer(
		name string,
		image container.Image,
		mounts ...container.Mount,
	) *container.Container
	GetContainer(name string) (*container.Container, error)
	RunContainer(con *container.Container) error
	StopContainer(con *container.Container) error
	RemoveContainer(con *container.Container) error
	AddToNetwork(n *container.Network, con *container.Container) error
}
