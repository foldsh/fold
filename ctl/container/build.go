// Package container contains utilities for building image from services.
package container

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

var (
	FailedToListImages   = errors.New("failed to list images")
	FailedToBuildImage   = errors.New("failed to build the image")
	FailedToPullImage    = errors.New("failed to pull the image")
	FailedToInspectImage = errors.New("failed to inspect the image")
)

type Image struct {
	ID      string
	Src     string
	Name    string
	WorkDir string
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

func (cr *ContainerRuntime) PullImage(name string) (*Image, error) {
	cr.logger.Infof("Pulling image %s", name)
	rc, err := cr.cli.ImagePull(cr.ctx, name, types.ImagePullOptions{})
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
	cr.logger.Infof("Successfully pulled image %s", name)
	img, err := cr.GetImage(name)
	if err != nil {
		return nil, err
	}
	// We're assuming that it will always be found as we just built it.
	return img, nil
}

func (cr *ContainerRuntime) ListImages() ([]*Image, error) {
	images, err := cr.cli.ImageList(cr.ctx, types.ImageListOptions{})
	if err != nil {
		return nil, FailedToListImages
	}
	var results []*Image
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if strings.HasPrefix(tag, "fold") {
				img := &Image{ID: img.ID, Name: tag}
				if err := cr.populateImageDetail(img); err != nil {
					return nil, err
				}
				results = append(results, img)
				break
			}
		}
	}
	return results, nil
}

func (cr *ContainerRuntime) GetImage(name string) (*Image, error) {
	// Woefully inefficient but simple and the API is shit for this use case.
	// Could look into caching ListImages if needs be.
	imgs, err := cr.ListImages()
	if err != nil {
		return nil, err
	}
	for _, img := range imgs {
		if img.Name == name {
			cr.logger.Debugf("found image %+v", img)
			return img, nil
		}
	}
	return nil, nil
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
	cr.logger.Debugf("building image %+v", img)
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
	builtImage, err := cr.GetImage(img.Name)
	if err != nil {
		return err
	}
	img.ID = builtImage.ID
	img.WorkDir = builtImage.WorkDir
	cr.logger.Debugf("finished building image %v", img)
	return nil
}

func (cr *ContainerRuntime) populateImageDetail(img *Image) error {
	detail, _, err := cr.cli.ImageInspectWithRaw(cr.ctx, img.ID)
	if err != nil {
		return err
	}
	img.ID = detail.ID
	img.WorkDir = detail.Config.WorkingDir
	return nil
}
