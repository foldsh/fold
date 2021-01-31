package build

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/foldsh/fold/logging"
)

type dockerClient struct {
	dc     *client.Client
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
		dc:     client,
		ctx:    ctx,
		logger: logger,
		out:    out,
	}, nil
}

func (dc *dockerClient) buildImage(image string, archive *archive.TempArchive) error {
	dc.logger.Debugf("building image")
	opts := types.ImageBuildOptions{
		Tags: []string{image},
	}

	resp, err := dc.dc.ImageBuild(dc.ctx, archive.File, opts)
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
