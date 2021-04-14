// This file contains a number of 'actions' which are common to many commands.
// It avoids duplicating things like common bits of resource acquisition across
// commands, and ensures that error handling is consistent for all of them.
// This also serves to make the commands much less verbose. When there are errors
// the only course of action is to stop the command with an appropriate help message.
// We can therefore capture a lot of the errors that library code throws in here.
package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
)

func projectHome(ctx *ctl.CmdCtx) string {
	s, err := project.Home()
	if err != nil {
		ctx.InformError(err)
		os.Exit(1)
	}
	return s
}

func loadProject(ctx *ctl.CmdCtx) *project.Project {
	p, err := project.Load(ctx, projectHome(ctx))
	if err != nil {
		if errors.Is(err, project.NotAFoldProject) {
			ctx.Inform(output.Error("this is not a fold project root."))
			ctx.Inform(output.Line("Please either initialise a project or cd to a project root."))
			os.Exit(1)
		} else if errors.Is(err, project.InvalidConfig) {
			ctx.Inform(output.Error("fold config is invalid."))
			ctx.Inform(output.Line("Please check that the yaml is valid and that you have spelled all the keys correctly."))
			os.Exit(1)
		} else {
			ctx.Inform(output.Error("failed to load fold config."))
			ctx.Inform(output.Line("Please ensure you're in a fold project root."))
			os.Exit(1)
		}
	}
	err = p.Validate()
	if err != nil {
		ctx.InformError(err)
	}
	return p
}

func loadProjectWithRuntime(ctx *ctl.CmdCtx, out io.Writer) *project.Project {
	dockerTimeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	p := loadProject(ctx)
	dc, err := container.NewDockerClient(ctx.Logger)
	if err != nil {
		ctx.InformError(err)
	}
	rt, err := container.NewRuntime(
		dockerTimeoutCtx,
		ctx.Logger,
		out,
		container.NewOSFileSystem(),
		dc,
	)
	if err != nil {
		ctx.InformError(err)
	}
	p.ConfigureContainerAPI(rt)
	return p
}

func saveProjectConfig(ctx *ctl.CmdCtx, p *project.Project) {
	err := p.SaveConfig(projectHome(ctx))
	if err != nil {
		ctx.Inform(output.Error("failed to solve fold config"))
		ctx.Inform(
			output.Line("Please check you have permission to write files in this directory."),
		)
		os.Exit(1)
	}
}

func getService(ctx *ctl.CmdCtx, p *project.Project, path string) *project.Service {
	service, err := p.GetService(path)
	if err != nil {
		ctx.Inform(output.Error(fmt.Sprintf("The path %s is not a registered service.", path)))
		ctx.Inform(
			output.Line(
				"Please check the path you typed or, if this is a mistake, make sure that the service",
			),
		)
		ctx.Inform(output.Line("is registered in your fold.yaml file."))
	}
	return service
}
