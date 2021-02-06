// Package container contains utilities for building image from services.
package container

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/foldsh/fold/logging"
	"github.com/moby/term"
)

type Image struct {
	Src  string
	Name string
}

type ImageBuilder interface {
	Build(spec Image) error
}

func NewImageBuilder(
	ctx context.Context, logger logging.Logger, out io.Writer,
) (ImageBuilder, error) {
	client, err := newDockerClient(logger)
	if err != nil {
		return nil, err
	}
	return &imageBuilder{client, ctx, logger, out}, nil
}

type imageBuilder struct {
	client DockerClient
	ctx    context.Context
	logger logging.Logger
	out    io.Writer
}

func (ib *imageBuilder) Build(img Image) error {
	// Create the workspace
	workDir, err := ioutil.TempDir("", "fold-build")
	if err != nil {
		ib.logger.Debugf("failed to create the temporary workspace: %v", err)
		return FailedToPrepareBuildArchive
	}
	defer os.RemoveAll(workDir)

	// Build the archive.
	archive, err := prepareArchive(img.Src, workDir, img.ignoreFilePatterns())
	if err != nil {
		ib.logger.Debugf("failed to build the tar archive for the build: %v", err)
		return FailedToPrepareBuildArchive
	}

	// Build the image
	ib.logger.Debugf("building image")
	opts := types.ImageBuildOptions{
		Tags: []string{img.Name},
	}

	resp, err := ib.client.ImageBuild(ib.ctx, archive.File, opts)
	if err != nil {
		ib.logger.Debugf("failed to build image %v", err)
		return FailedToBuildImage
	}
	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)

	if err = jsonmessage.DisplayJSONMessagesStream(
		resp.Body, ib.out, termFd, isTerm, nil,
	); err != nil {
		ib.logger.Debugf("failed to build image %v", err)
		return FailedToBuildImage
	}
	return nil
}

func (spec *Image) ignoreFilePatterns() []string {
	gitIgnore := filepath.Join(spec.Src, ".gitignore")
	if _, err := os.Stat(gitIgnore); os.IsNotExist(err) {
		return []string{}
	}

	bytes, err := ioutil.ReadFile(gitIgnore)
	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(string(bytes), "\n")
}
