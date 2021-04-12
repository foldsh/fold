package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/ctl/git"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/version"
)

func NewNewCmd(ctx *ctl.CmdCtx) *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new [resource]",
		Short: "Create new fold resources",
		Long:  "Create new fold resources",
	}
	newCmd.AddCommand(NewProjectCommand(ctx))
	newCmd.AddCommand(NewServiceCommand(ctx))
	return newCmd
}

var projectLong = trimf(
	`Create a new fold project in the current directory or in the specified path.

If you do not specify a path then the project is created in the current directory and you will 
be prompted for a project name.

If you specify a path then the end of the path will be selected as the project name. 
For example, running:

foldctl new project fold/my-project

Will result in a project with the name 'my-project' located in the directory fold. All directories
in the given path will be created if they do not exist already.

Project names are validated against the regex %s.

Here are some examples of valid project names:
- FoldProject
- fold-project
- fold_project
- Fold-Project

Some examples of invalid project names are:
- _FoldProject
- 123_FoldProject
- Fold/Project

Aside from the name, you will be prompted for other values required by the fold config file.`,
	project.ProjectNameRegex,
)

func NewProjectCommand(ctx *ctl.CmdCtx) *cobra.Command {
	return &cobra.Command{
		Use:     "project [path]",
		Example: "foldctl new project\nfoldctl new project path/to/new-project",
		Short:   "Create a new fold project",
		Long:    projectLong,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				projectPath string
				projectName string
				mkDir       bool // Ugly but overall nicer than the alternative imo.

				projectNameValidator = newRegexValidator(ctx, project.ProjectNameRegex, "project")
			)
			if len(args) != 0 {
				projectPath = args[0]
				abs, err := filepath.Abs(projectPath)
				if err != nil {
					ctx.Inform(output.Error(fmt.Sprintf("%s is not a valid path", projectPath)))
				}
				projectName = filepath.Base(abs)
				err = projectNameValidator(projectName)
				if err != nil {
					ctx.Inform(output.Error(err.Error()))
					os.Exit(1)
				}
				mkDir = true
			} else {
				projectPath = "."
			}
			// If the project path is already a fold project we bail.
			if project.IsAFoldProject(projectPath) {
				ctx.Inform(output.Error(fmt.Sprintf("%s is already a fold project.", projectPath)))
				os.Exit(1)
			}

			// We're all clear so lets prompt for the project name.
			prompt := fmt.Sprintf("Name (must match %s)", project.ProjectNameRegex)
			projectName = runPrompt(ctx, prompt, projectNameValidator)

			// Create the project directory if we need to
			if mkDir {
				if err := os.MkdirAll(projectPath, fs.DIR_PERMISSIONS); err != nil {
					ctx.Inform(
						output.Error(
							fmt.Sprintf(
								"failed to create the project directory at %s",
								projectPath,
							),
						),
					)
					os.Exit(1)
				}
			}
			// Prompt for the rest of the details
			maintainer := runPrompt(ctx, "Maintainer", output.NoopValidator)
			email := runPrompt(ctx, "Email", output.NoopValidator)
			repo := runPrompt(ctx, "Repository", output.NoopValidator)

			p := &project.Project{
				Name:       projectName,
				Maintainer: maintainer,
				Email:      email,
				Repository: repo,
			}
			p.ConfigureCmdCtx(ctx)
			saveProjectConfig(ctx, p)
			ctx.Inform(
				output.Success(fmt.Sprintf("Successfully created the project %s", p.Name)),
			)
		},
	}
}

var serviceLong = trimf(`
Creates a new fold service from a template.

The command will run you through a series of prompts to create your new service. You will be
required to choose a project name, a template and a language. The prompts will give you the
available options to choose from for the templates and languages.

If you want to browse the available templates head to https://github.com/foldsh/templates.

The service name will be validated against the regex %s and the prompt will indicate when
the name you have entered is valid.

This command can only be run from a fold project root and it will create the service
relative to the current project root. The service will be created in a directory with the same
name as the service itself.
`, project.ServiceNameRegex)

func NewServiceCommand(ctx *ctl.CmdCtx) *cobra.Command {
	return &cobra.Command{
		Use:     "service",
		Example: "foldctl new service",
		Short:   "Create a new fold service",
		Long:    serviceLong,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// First up try to load the project to make sure we're in a project root.
			p := loadProject(ctx)
			// We're good to go, lets update the templates
			updateTemplates(ctx)
			// Ok lets prompt for the service name.

			serviceNameValidator := newRegexValidator(ctx, project.ServiceNameRegex, "service")
			namePrompt := fmt.Sprintf("Name (must match %s)", project.ServiceNameRegex)
			name := runPrompt(ctx, namePrompt, serviceNameValidator)
			// TODO we can generate the list of template and language options dynamically but this is
			// fine for now.
			// And now the template
			template := runSelect(ctx, "Template", []string{"basic"})
			// And finally the language
			language := runSelect(ctx, "Language", []string{"go", "js", "ts"})
			// Build the absolute path to the new service.
			servicePath := filepath.Join(".", name)
			absPath, err := filepath.Abs(servicePath)
			if err != nil {
				ctx.Inform(servicePathInvalid)
			}

			// Check if the directory is empty
			empty, err := fs.IsEmpty(absPath)
			if err == nil && !empty {
				ctx.Inform(
					output.Error(
						fmt.Sprintf(
							"the target directory %s already exists and is not empty.",
							absPath,
						),
					),
				)
				ctx.Inform(
					output.Line(
						"Please either choose a different name for your service or remove the existing directory.",
					),
				)
				os.Exit(1)
			}

			// And create the path to the relevant template
			templatePath := filepath.Join(ctx.Config.FoldTemplates, template, language)

			// Create the directory for the new service.
			ctx.Logger.Debugf("Creating service directory")
			err = os.MkdirAll(absPath, fs.DIR_PERMISSIONS)
			if err != nil {
				ctx.Inform(output.Error("Failed to create a directory at the path you specified."))
				ctx.Inform(output.Line(checkPermissions))
			}

			// Copy the contents of the chosen template into the new directory.
			ctx.Logger.Debugf("Copying project template to service directory")
			err = fs.CopyDir(templatePath, absPath)
			if err != nil {
				os.RemoveAll(absPath)
				ctx.Inform(
					output.Error("failed to create a project template at the path you specified."),
				)
				ctx.Inform(output.Line(checkPermissions))
				os.Exit(1)
			}

			// We're just about done, create the service so we can add it to the project
			service := project.Service{
				Name: name,
				Path: servicePath,
			}

			// Update the project config
			ctx.Logger.Debugf("Adding service to project")
			p.AddService(service)
			ctx.Logger.Debugf("Saving project config")
			saveProjectConfig(ctx, p)
			ctx.Inform(
				output.Success(fmt.Sprintf("Successfully created the service %s", service.Name)),
			)
		},
	}
}

func runSelect(ctx *ctl.CmdCtx, label string, items []string) string {
	value, err := ctx.Select(label, items)
	if err != nil {
		ctx.Inform(output.Error(err.Error()))
		os.Exit(1)
	}
	return value
}

func runPrompt(ctx *ctl.CmdCtx, label string, validator func(string) error) string {
	value, err := ctx.Prompt(label, validator)
	if err != nil {
		ctx.Inform(output.Error(err.Error()))
		os.Exit(1)
	}
	return value
}

func evalPrompt(ctx *ctl.CmdCtx, value, label string, err error) string {
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			os.Exit(1)
		} else {
			ctx.Inform(output.Error(fmt.Sprintf("specified %s is not valid.", label)))
			os.Exit(1)
		}
	}
	return value
}

func newRegexValidator(ctx *ctl.CmdCtx, regex, target string) output.Validator {
	return func(input string) error {
		err := output.RegexValidator(regex, input)
		if err != nil {
			if errors.Is(err, output.RegexMatchError) {
				return fmt.Errorf(
					"%s is not a valid %s name. It must match the regex %s",
					input,
					target,
					regex,
				)
			}
			return fmt.Errorf("failed to validate project name %s", input)
		}
		return nil
	}
}

func updateTemplates(ctx *ctl.CmdCtx) {
	// Get or update the templates.
	ctx.Inform(output.Line("Updating the templates repository..."))
	out := ctx.InformWriter(output.WithPrefix(output.Blue("git: ")))
	err := git.UpdateTemplates(out, ctx.Config.FoldTemplates, version.FoldVersion.String())
	if err != nil {
		ctx.Inform(output.Error(err.Error()))
		os.Exit(1)
	}
}
