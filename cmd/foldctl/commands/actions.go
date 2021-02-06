// This file contains a number of 'actions' which are common to many commands.
// It avoids duplicating things like common bits of resource acquisition across
// commands, and ensures that error handling is consistent for all of them.
// This also serves to make the commands much less verbose. When there are errors
// the only course of action is to stop the command with an appropriate help message.
// We can therefore capture a lot of the errors that library code throws in here.
package commands

import (
	"errors"
	"fmt"

	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/project"
)

func loadProject() *project.Project {
	p, err := project.Load(logger)
	if err != nil {
		if errors.Is(err, project.NotAFoldProject) {
			exitWithMessage(
				"This is not a fold project root.",
				"Please either initialise a project or cd to a project root.",
			)
		} else if errors.Is(err, project.InvalidConfig) {
			exitWithMessage(
				"Fold config is invalid.",
				"Please check that the yaml is valid and that you have spelled all the keys correctly.",
			)
		} else {
			exitWithMessage("Failed to load fold config. Please ensure you're in a fold project root.")
		}
	}
	return p
}

func newOut(outPrefix string) *streamLinePrefixer {
	return newStreamLinePrefixer(serr, blue(outPrefix))
}

func loadProjectWithRuntime(out *streamLinePrefixer) *project.Project {
	p := loadProject()
	rt, err := container.NewRuntime(commandCtx, logger, out)
	exitIfError(err, "Failed to create DockerClient. Ensure that the docker daemon is running.")
	p.ConfigureBackend(rt)
	return p
}

func saveProjectConfig(p *project.Project) {
	err := p.SaveConfig()
	exitIfError(
		err,
		"Failed to save fold config.",
		"Please check you have permission to write files in this directory.",
	)
}

func getService(p *project.Project, path string) *project.Service {
	service, err := p.GetService(path)
	exitIfError(
		err,
		fmt.Sprintf("The path %s is not a registered service.", path),
		"Please check the path you typed or, if this is a mistake, make sure that the service",
		"is registered in your fold.yaml file.",
	)
	return service
}
