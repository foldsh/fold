package container

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

var (
	DockerEngineError                  = errors.New("failed to connect to docker engine")
	FailedToConnectToDockerEngineError = errors.New("failed to connect to docker engine")
	FailedToPrepareBuildArchive        = errors.New("failed to prepare the build archive")
	FailedToCreateContainer            = errors.New("failed to create the container")
	FailedToStartContainer             = errors.New("failed to start the container")
	FailedToStopContainer              = errors.New("failed to stop the container")
	FailedToRemoveContainer            = errors.New("failed to remove the container")

	foldPrefix = "fold."
)

// A very simple struct for representing the functionality we actually
// need to expose from here.
type Container struct {
	ID           string
	Name         string
	NetworkAlias string
	Image        Image
	Ports        []int
	Mounts       []Mount
}

type Mount struct {
	Src string
	Dst string
}

func (cr *ContainerRuntime) NewContainer(
	name string, image Image, mounts ...Mount,
) *Container {
	return &Container{
		Name:   fmt.Sprintf("%s%s", foldPrefix, name),
		Image:  image,
		Mounts: mounts,
	}
}

func (cr *ContainerRuntime) AllContainers() ([]*Container, error) {
	return cr.listContainers()
}

func (cr *ContainerRuntime) GetContainer(name string) (*Container, error) {
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

func (cr *ContainerRuntime) RunContainer(net *Network, con *Container) error {
	cr.logger.Debugf("Building container %v in network %v", con, net)
	portBindings := map[nat.Port][]nat.PortBinding{}
	for _, p := range con.Ports {
		binding := []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", p)},
		}
		portBindings[nat.Port("6123/tcp")] = binding
	}
	var mounts []mount.Mount
	for _, m := range con.Mounts {
		mounts = append(mounts, mount.Mount{
			// Type:   mount.TypeVolume,
			Type:   mount.TypeBind,
			Source: m.Src,
			Target: m.Dst,
		})
	}
	var watchDir string
	if len(mounts) > 0 {
		watchDir = mounts[0].Target
	}
	resp, err := cr.cli.ContainerCreate(
		cr.ctx,
		&container.Config{
			Image: con.Image.Name,
			Env: []string{
				"FOLD_STAGE=LOCAL",
				fmt.Sprintf("FOLD_WATCH_DIR=%s", watchDir),
			},
		},
		// TODO make auto removing containers configurable
		&container.HostConfig{PortBindings: portBindings, Mounts: mounts, AutoRemove: true},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				net.Name: &network.EndpointSettings{Aliases: []string{con.NetworkAlias}},
			},
		},
		nil,
		con.Name,
	)
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

func (cr *ContainerRuntime) StopContainer(con *Container) error {
	if err := cr.cli.ContainerStop(cr.ctx, con.ID, nil); err != nil {
		return FailedToStopContainer
	}
	return nil
}

func (cr *ContainerRuntime) RemoveContainer(con *Container) error {
	if err := cr.cli.ContainerRemove(cr.ctx, con.ID, types.ContainerRemoveOptions{}); err != nil {
		cr.logger.Debugf("Failed to remove the container %s: %v", con.Name, err)
		return FailedToRemoveContainer
	}
	return nil
}

func (cr *ContainerRuntime) AddToNetwork(n *Network, con *Container) error {
	err := cr.cli.NetworkConnect(
		cr.ctx,
		n.ID,
		con.ID,
		&network.EndpointSettings{Aliases: []string{con.NetworkAlias}},
	)
	if err != nil {
		return FailedToJoinNetwork
	}
	return nil
}

func (cr *ContainerRuntime) RemoveFromNetwork(net *Network, con *Container) error {
	err := cr.cli.NetworkDisconnect(cr.ctx, net.ID, con.ID, false)
	if err != nil {
		return FailedToLeaveNetwork
	}
	return nil
}

func (cr *ContainerRuntime) listContainers() ([]*Container, error) {
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
				var mounts []Mount
				for _, mp := range c.Mounts {
					mounts = append(mounts, Mount{mp.Source, mp.Destination})
				}
				foldContainers = append(
					foldContainers,
					&Container{ID: c.ID, Name: name, Image: Image{Name: c.Image}, Mounts: mounts},
				)
			}
		}
	}
	return foldContainers, nil
}
