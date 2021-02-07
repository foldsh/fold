// Package container contains utilities for building image from services.
package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

type Image struct {
	Src  string
	Name string
}

func (cr *ContainerRuntime) PullImage(img string) (*Image, error) {
	cr.logger.Infof("Pulling image %s", img)
	rc, err := cr.cli.ImagePull(cr.ctx, img, types.ImagePullOptions{})
	if err != nil {
		return nil, FailedToPullImage
	}
	defer rc.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)

	if err = jsonmessage.DisplayJSONMessagesStream(
		rc, cr.out, termFd, isTerm, nil,
	); err != nil {
		cr.logger.Debugf("failed to pull image %v", err)
		return nil, FailedToPullImage
	}
	cr.logger.Infof("Successfully pulled image %s", img)
	return &Image{Name: img}, nil
}

func (cr *ContainerRuntime) BuildImage(img *Image) error {
	// Create the workspace
	workDir, err := ioutil.TempDir("", "fold-build")
	if err != nil {
		cr.logger.Debugf("failed to create the temporary workspace: %v", err)
		return FailedToPrepareBuildArchive
	}
	defer os.RemoveAll(workDir)

	// Build the archive.
	archive, err := prepareArchive(img.Src, workDir, img.ignoreFilePatterns())
	if err != nil {
		cr.logger.Debugf("failed to build the tar archive for the build: %v", err)
		return FailedToPrepareBuildArchive
	}

	// Build the image
	cr.logger.Debugf("building image")
	opts := types.ImageBuildOptions{
		Tags: []string{img.Name},
	}

	resp, err := cr.cli.ImageBuild(cr.ctx, archive.File, opts)
	if err != nil {
		cr.logger.Debugf("failed to build image %v", err)
		return FailedToBuildImage
	}
	defer resp.Body.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)

	if err = jsonmessage.DisplayJSONMessagesStream(
		resp.Body, cr.out, termFd, isTerm, nil,
	); err != nil {
		cr.logger.Debugf("failed to build image %v", err)
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
