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
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/foldsh/fold/logging"
	"github.com/moby/term"
)

type ImageSpec struct {
	Src          string
	Name         string
	Logger       logging.Logger
	Out          io.Writer
	DockerClient DockerClient

	workDir string
}

func (spec *ImageSpec) Build(ctx context.Context) error {
	// Create the workspace
	workDir, err := ioutil.TempDir("", "fold-build")
	if err != nil {
		spec.Logger.Debugf("failed to create the temporary workspace: %v", err)
		return FailedToPrepareBuildArchive
	}
	spec.workDir = workDir
	defer os.RemoveAll(workDir)

	// Build the archive.
	archive, err := prepareArchive(spec.Src, workDir, spec.ignoreFilePatterns())
	if err != nil {
		spec.Logger.Debugf("failed to build the tar archive for the build: %v", err)
		return FailedToPrepareBuildArchive
	}

	// Build the image
	err = buildImage(ctx, spec, archive)
	if err != nil {
		spec.Logger.Debugf("failed to build image %v", err)
		return FailedToBuildImage
	}
	return nil
}

func (bs *ImageSpec) ignoreFilePatterns() []string {
	gitIgnore := filepath.Join(bs.Src, ".gitignore")
	if _, err := os.Stat(gitIgnore); os.IsNotExist(err) {
		return []string{}
	}

	bytes, err := ioutil.ReadFile(gitIgnore)
	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(string(bytes), "\n")
}

func buildImage(
	ctx context.Context, spec *ImageSpec, archive *archive.TempArchive,
) error {
	spec.Logger.Debugf("building image")
	opts := types.ImageBuildOptions{
		Tags: []string{spec.Name},
	}

	resp, err := spec.DockerClient.ImageBuild(ctx, archive.File, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)

	if err := jsonmessage.DisplayJSONMessagesStream(resp.Body, spec.Out, termFd, isTerm, nil); err != nil {
		return err
	}
	return nil
}
