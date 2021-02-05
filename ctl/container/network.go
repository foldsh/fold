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

	rt *containerRuntime
}

func (net *Network) CreateIfNotExists() error {
	exists, err := net.Exists()
	if err != nil {
		return err
	} else if exists {
		return nil
	} else {
		return net.rt.createNetwork(net)
	}
}

func (net *Network) Remove() error {
	return net.rt.removeNetwork(net)
}

func (net *Network) RemoveIfExists() error {
	exists, err := net.Exists()
	if err != nil {
		return err
	} else if !exists {
		return nil
	} else {
		return net.rt.removeNetwork(net)
	}
}

func (net *Network) Exists() (bool, error) {
	// Best just to always look it up rather than checking if it's
	// already on the Network.
	existingNetworks, err := net.rt.cli.NetworkList(net.rt.ctx, types.NetworkListOptions{})
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
