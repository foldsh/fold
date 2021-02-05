package container

import (
	"strings"
	"unicode/utf8"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/foldsh/fold/logging"
)

func NewDockerClient(logger logging.Logger) (DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Debugf("failed to initialised docker client")
		return nil, FailedToConnectToDockerEngineError
	}
	return client, nil
}

func (cr *containerRuntime) createNetwork(net *Network) error {
	networkRes, err := cr.cli.NetworkCreate(cr.ctx, net.Name, types.NetworkCreate{})
	if err != nil {
		return FailedToCreateNetwork
	}
	net.ID = networkRes.ID
	return nil
}

func (cr *containerRuntime) removeNetwork(net *Network) error {
	err := cr.cli.NetworkRemove(cr.ctx, net.ID)
	if err != nil {
		return FailedToDestroyNetwork
	}
	return nil
}

func (cr *containerRuntime) addToNetwork(n *Network, con *Container) error {
	err := cr.cli.NetworkConnect(cr.ctx, n.ID, con.ID, &network.EndpointSettings{})
	if err != nil {
		return FailedToJoinNetwork
	}
	return nil
}

func (cr *containerRuntime) removeFromNetwork(net *Network, con *Container) error {
	err := cr.cli.NetworkDisconnect(cr.ctx, net.ID, con.ID, false)
	if err != nil {
		return FailedToLeaveNetwork
	}
	return nil
}

func (cr *containerRuntime) pullImage(image string) error {
	// TODO returns a ReadCloser - pipe this through to the cli perhaps?
	_, err := cr.cli.ImagePull(cr.ctx, image, types.ImagePullOptions{})
	if err != nil {
		return FailedToPullImage
	}
	return nil
}

func (cr *containerRuntime) runContainer(con *Container) error {
	resp, err := cr.cli.ContainerCreate(cr.ctx, &container.Config{
		Image: con.Image,
	}, nil, nil, nil, con.Name)
	if err != nil {
		cr.logger.Debugf("Failed to create container: %v", err)
		return FailedToCreateContainer
	}

	if err := cr.cli.ContainerStart(cr.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return FailedToStartContainer
	}
	con.ID = resp.ID
	return nil
}

func (cr *containerRuntime) stopContainer(con *Container) error {
	if err := cr.cli.ContainerStop(cr.ctx, con.ID, nil); err != nil {
		return FailedToStopContainer
	}
	return nil
}

func (cr *containerRuntime) removeContainer(con *Container) error {
	if err := cr.cli.ContainerRemove(cr.ctx, con.ID, types.ContainerRemoveOptions{}); err != nil {
		return FailedToRemoveContainer
	}
	return nil
}

func (cr *containerRuntime) listContainers() ([]*Container, error) {
	containers, err := cr.cli.ContainerList(cr.ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, DockerEngineError
	}
	var foldContainers []*Container
	for _, c := range containers {
		cr.logger.Debugf("Examining container %s with names %v", c.ID, c.Names)
		for _, name := range c.Names {
			// Internally, docker represents container names with a leading slash, so we remove that
			// to make sure we spot the fold containers. It does this because it means that the
			// container is running on the local engine.
			_, i := utf8.DecodeRuneInString(name)
			name := name[i:]
			if strings.HasPrefix(name, foldPrefix) {
				cr.logger.Debugf("Identified container with name %s as fold container", name)
				var volumes []Volume
				for _, mp := range c.Mounts {
					volumes = append(volumes, Volume{mp.Source, mp.Destination})
				}
				foldContainers = append(
					foldContainers,
					&Container{ID: c.ID, Name: name, Image: c.Image, Volumes: volumes, rt: cr},
				)
			}
		}
	}
	return foldContainers, nil
}
