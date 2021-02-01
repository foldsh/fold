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

	"github.com/foldsh/fold/logging"
)

type BuildSpec struct {
	Src    string
	Image  string
	Logger logging.Logger
	Out    io.Writer

	workDir string
}

func Build(ctx context.Context, spec *BuildSpec) error {
	// Create this first as if it fails there is no point building the archive
	dc, err := newDockerClient(ctx, spec.Logger, spec.Out)
	if err != nil {
		return err
	}
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
	err = dc.buildImage(spec.Image, archive)
	if err != nil {
		spec.Logger.Debugf("failed to build image %v", err)
		return FailedToBuildImage
	}
	return nil
}

func (bs *BuildSpec) ignoreFilePatterns() []string {
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
