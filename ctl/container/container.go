package container

import "errors"

var (
	FailedToConnectToDockerEngineError      = errors.New("failed to connect to docker engine")
	FailedToDetermineDockerEngineAPIVersion = errors.New("failed to set the docker engine API version")
	FailedToCreateNetwork                   = errors.New("failed to create the network")
	FailedToDestroyNetwork                  = errors.New("failed to destroy the network")
	FailedToJoinNetwork                     = errors.New("failed to join the network")
	FailedToLeaveNetwork                    = errors.New("failed to leave the network")
	FailedToPrepareBuildArchive             = errors.New("failed to prepare the build archive")
	FailedToBuildImage                      = errors.New("failed to build the image")
	FailedToPullImage                       = errors.New("failed to pull the image")
	FailedToCreateContainer                 = errors.New("failed to create the container")
	FailedToStartContainer                  = errors.New("failed to start the container")
	FailedToStopContainer                   = errors.New("failed to stop the container")
)

// A very simple struct for representing the functionality we actually need to expose
// from here.
type Container struct {
	Name    string
	Image   string
	Volumes string

	id string
}

func (c *Container) Stop() error {
	return nil
}

func (c *Container) JoinNetwork(net *Network) error {
	return nil
}

func (c *Container) LeaveNetwork(net *Network) error {
	return nil
}

// A very simple abstraction for a network that avoids exposing too much of the docker
// API
type Network struct {
	Name       string
	Containers []Container

	id string
}

func (c *Network) Shutdown() error {
	return nil
}
