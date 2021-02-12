package container

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerClient interface {
	ImagePull(
		ctx context.Context,
		ref string,
		options types.ImagePullOptions,
	) (io.ReadCloser, error)

	ImageList(
		ctx context.Context,
		options types.ImageListOptions,
	) ([]types.ImageSummary, error)

	ImageBuild(
		ctx context.Context,
		buildContext io.Reader,
		options types.ImageBuildOptions,
	) (types.ImageBuildResponse, error)

	ImageInspectWithRaw(
		ctx context.Context,
		imageID string,
	) (types.ImageInspect, []byte, error)

	NetworkCreate(
		ctx context.Context,
		name string,
		options types.NetworkCreate,
	) (types.NetworkCreateResponse, error)

	NetworkRemove(
		ctx context.Context,
		networkID string,
	) error

	NetworkList(
		ctx context.Context,
		options types.NetworkListOptions,
	) ([]types.NetworkResource, error)

	NetworkConnect(
		ctx context.Context,
		networkID,
		containerID string,
		config *network.EndpointSettings,
	) error

	NetworkDisconnect(
		ctx context.Context,
		networkID,
		containerID string,
		force bool,
	) error

	ContainerCreate(
		ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		platform *specs.Platform,
		containerName string,
	) (container.ContainerCreateCreatedBody, error)

	ContainerStart(
		ctx context.Context,
		containerID string,
		options types.ContainerStartOptions,
	) error

	ContainerStop(
		ctx context.Context,
		containerID string,
		timeout *time.Duration,
	) error

	ContainerRemove(
		ctx context.Context,
		containerID string,
		options types.ContainerRemoveOptions,
	) error

	ContainerList(
		ctx context.Context,
		options types.ContainerListOptions,
	) ([]types.Container, error)
}
