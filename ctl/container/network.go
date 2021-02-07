package container

import (
	"errors"

	"github.com/docker/docker/api/types"
)

var (
	FailedToCreateNetwork  = errors.New("failed to create the network")
	FailedToDestroyNetwork = errors.New("failed to destroy the network")
	FailedToJoinNetwork    = errors.New("failed to join the network")
	FailedToLeaveNetwork   = errors.New("failed to leave the network")
)

// A very simple abstraction for a network that avoids exposing too much of the docker
// API
type Network struct {
	ID   string
	Name string
}

func (cr *ContainerRuntime) NewNetwork(name string) *Network {
	return &Network{Name: name}
}

func (cr *ContainerRuntime) CreateNetwork(net *Network) error {
	networkRes, err := cr.cli.NetworkCreate(cr.ctx, net.Name, types.NetworkCreate{})
	if err != nil {
		return FailedToCreateNetwork
	}
	net.ID = networkRes.ID
	return nil
}

func (cr *ContainerRuntime) RemoveNetwork(net *Network) error {
	err := cr.cli.NetworkRemove(cr.ctx, net.ID)
	if err != nil {
		return FailedToDestroyNetwork
	}
	return nil
}

func (cr *ContainerRuntime) NetworkExists(net *Network) (bool, error) {
	// Best just to always look it up rather than checking if it's
	// already on the Network.
	existingNetworks, err := cr.cli.NetworkList(cr.ctx, types.NetworkListOptions{})
	if err != nil {
		return false, DockerEngineError
	}
	for _, existing := range existingNetworks {
		if existing.Name == net.Name {
			net.ID = existing.ID
			return true, nil
		}
	}
	return false, nil
}
