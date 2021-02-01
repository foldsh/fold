package container

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/foldsh/fold/logging"
)

type dockerClient struct {
	cli    *client.Client
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
}

func newDockerClient(ctx context.Context, logger logging.Logger, out io.Writer) (*dockerClient, error) {
	apiVersion, err := determineDockerAPIVersion(ctx, logger)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("DOCKER_API_VERSION", apiVersion)
	if err != nil {
		return nil, FailedToDetermineDockerEngineAPIVersion
	}
	client, err := client.NewEnvClient()
	if err != nil {
		logger.Debugf("failed to initialised docker client")
		return nil, FailedToConnectToDockerEngineError
	}
	return &dockerClient{
		cli:    client,
		ctx:    ctx,
		logger: logger,
		out:    out,
	}, nil
}

// In order to support a wide range of versions of the docker engine, we need to find
// out what version of the API the current engine supports. The only way to do this is to
// create a client, ping the engine, and extract the APIVersion from the response.
// It means we end up creating a client and throwing it away just for this but oh well.
func determineDockerAPIVersion(ctx context.Context, logger logging.Logger) (string, error) {
	logger.Debugf("determining supported docker engine api version")
	client, err := client.NewEnvClient()
	if err != nil {
		logger.Debugf("failed to initialised docker client: %v", err)
		return "", FailedToConnectToDockerEngineError
	}
	ping, err := client.Ping(ctx)
	if err != nil {
		logger.Debugf("failed to ping docker client: %v", err)
		return "", FailedToConnectToDockerEngineError
	}
	version := ping.APIVersion
	if version == "" {
		logger.Debugf("failed to determine supported docker engine api version")
		return "", FailedToConnectToDockerEngineError
	}
	logger.Debugf("determined supported docker engine api version to be %s", version)
	return version, nil
}

func (dc *dockerClient) buildImage(image string, archive *archive.TempArchive) error {
	dc.logger.Debugf("building image")
	opts := types.ImageBuildOptions{
		Tags: []string{image},
	}

	resp, err := dc.cli.ImageBuild(dc.ctx, archive.File, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)

	if err := jsonmessage.DisplayJSONMessagesStream(resp.Body, dc.out, termFd, isTerm, nil); err != nil {
		return err
	}
	return nil
}

func (dc *dockerClient) createNetwork(ctx context.Context, net *Network) error {
	networkRes, err := dc.cli.NetworkCreate(ctx, net.Name, types.NetworkCreate{})
	if err != nil {
		return FailedToCreateNetwork
	}
	net.id = networkRes.ID
	return nil
}

func (dc *dockerClient) destroyNetwork(ctx context.Context, net *Network) error {
	err := dc.cli.NetworkRemove(ctx, net.id)
	if err != nil {
		return FailedToDestroyNetwork
	}
	return nil
}

func (dc *dockerClient) addToNetwork(
	ctx context.Context, n *Network, con *Container,
) error {
	err := dc.cli.NetworkConnect(ctx, n.id, con.id, &network.EndpointSettings{})
	if err != nil {
		return FailedToJoinNetwork
	}
	return nil
}

func (dc *dockerClient) removeFromNetwork(
	ctx context.Context, net *Network, con *Container,
) error {
	err := dc.cli.NetworkDisconnect(ctx, net.id, con.id, false)
	if err != nil {
		return FailedToLeaveNetwork
	}
	return nil
}

func (dc *dockerClient) pullImage(ctx context.Context, image string) error {
	// TODO returns a ReadCloser - pipe this through to the cli perhaps?
	_, err := dc.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return FailedToPullImage
	}
	return nil
}

func (dc *dockerClient) runContainer(ctx context.Context, con *Container) error {
	resp, err := dc.cli.ContainerCreate(ctx, &container.Config{
		Image: con.Image,
	}, nil, nil, nil, con.Name)
	if err != nil {
		return FailedToCreateContainer
	}

	if err := dc.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return FailedToStartContainer
	}
	return nil
}

func (dc *dockerClient) stopContainer(ctx context.Context, con *Container) error {
	// TODO get the container ID for the passed name
	if err := dc.cli.ContainerStop(ctx, con.id, nil); err != nil {
		return FailedToStopContainer
	}
	return nil
}
