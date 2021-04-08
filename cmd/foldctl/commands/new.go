package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/ctl/git"
	"github.com/foldsh/fold/ctl/project"
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newProjectCmd)
	newCmd.AddCommand(newServiceCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [resource]",
	Short: "Create new fold resources",
	Long:  "Create new fold resources",
}

var newProjectCmd = &cobra.Command{
	Use:     "project [path]",
	Example: "foldctl new project\nfoldctl new project path/to/new-project",
	Short:   "Create a new fold project",
	Long: fmt.Sprintf(
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
	),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			projectPath string
			projectName string
			mkDir       bool // Ugly but overall nicer than the alternative imo.
		)
		if len(args) != 0 {
			projectPath = args[0]
			abs, err := filepath.Abs(projectPath)
			exitIfErr(err, fmt.Sprintf("%s is not a valid path", projectPath))
			projectName = filepath.Base(abs)
			err = projectNameValidator(projectName)
			exitIfErr(err)
			mkDir = true
		} else {
			projectPath = "."
			prompt := fmt.Sprintf("Name (must match %s)", project.ProjectNameRegex)
			projectName = runPrompt(promptui.Prompt{Label: prompt, Validate: projectNameValidator})
		}
		if project.IsAFoldProject(projectPath) {
			exitWithMessage(fmt.Sprintf("%s is already a fold project.", projectPath))
		}
		// Create the project directory if we need to
		if mkDir {
			if err := os.MkdirAll(projectPath, fs.DIR_PERMISSIONS); err != nil {
				exitWithMessage(
					fmt.Sprintf("Failed to create the project directory at %s", projectPath),
				)
			}
		}
		// Prompt for the rest of the details
		maintainer := runPrompt(promptui.Prompt{Label: "Maintainer"})
		email := runPrompt(promptui.Prompt{Label: "Email"})
		repo := runPrompt(promptui.Prompt{Label: "Repository"})

		p := &project.Project{
			Name:       projectName,
			Maintainer: maintainer,
			Email:      email,
			Repository: repo,
		}
		p.ConfigureLogger(logger)
		saveProjectConfig(p)
	},
}

var newServiceCmd = &cobra.Command{
	Use:     "service",
	Example: "foldctl new service",
	Short:   "Create a new fold service",
	Long: fmt.Sprintf(`Creates a new fold service from a template.

The command will run you through a series of prompts to create your new service. You will be
required to choose a project name, a template and a language. The prompts will give you the
available options to choose from for the templates and languages.

If you want to browse the available templates head to https://github.com/foldsh/templates.

The service name will be validated against the regex %s and the prompt will indicate when
the name you have entered is valid.

This command can only be run from a fold project root and it will create the service
relative to the current project root. The service will be created in a directory with the same
name as the service itself.
`, project.ServiceNameRegex),
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// First up try to load the project to make sure we're in a project root.
		p := loadProject()
		// We're good to go, lets update the templates
		updateTemplates()
		// Ok lets prompt for the service name.
		namePrompt := fmt.Sprintf("Name (must match %s)", project.ServiceNameRegex)
		name := runPrompt(promptui.Prompt{Label: namePrompt, Validate: serviceNameValidator})
		// TODO we can generate the list of template and language options dynamically but this is
		// fine for now.
		// And now the template
		template := runSelect(promptui.Select{Label: "Template", Items: []string{"basic"}})
		// And finally the language
		language := runSelect(
			promptui.Select{Label: "Language", Items: []string{"go", "js", "ts"}},
		)
		// Build the absolute path to the new service.
		servicePath := filepath.Join(".", name)
		absPath, err := filepath.Abs(servicePath)
		exitIfErr(err, servicePathInvalid)

		// Check if the directory is empty
		empty, err := fs.IsEmpty(absPath)
		if err == nil && !empty {
			exitWithMessage(
				fmt.Sprintf("The target directory %s already exists and is not empty.", absPath),
				"Please either choose a different name for your service or remove the existing directory.",
			)
		}

		// And create the path to the relevant template
		templatePath := filepath.Join(foldTemplates, template, language)

		// Create the directory for the new service.
		logger.Debugf("Creating service directory")
		err = os.MkdirAll(absPath, fs.DIR_PERMISSIONS)
		exitIfErr(err, "Failed to create a directory at the path you specified.", checkPermissions)

		// Copy the contents of the chosen template into the new directory.
		logger.Debugf("Copying project template to service directory")
		err = fs.CopyDir(templatePath, absPath)
		if err != nil {
			os.RemoveAll(absPath)
			exitWithMessage(
				"Failed to create a project template at the path you specified.",
				checkPermissions,
			)
		}

		// We're just about done, create the service so we can add it to the project
		service := project.Service{
			Name: name,
			Path: servicePath,
		}

		// Update the project config
		logger.Debugf("Adding service to project")
		p.AddService(service)
		logger.Debugf("Saving project config")
		saveProjectConfig(p)
	},
}

func runSelect(sel promptui.Select) string {
	_, value, err := sel.Run()
	return evalPrompt(value, sel.Label.(string), err)
}

func runPrompt(prompt promptui.Prompt) string {
	value, err := prompt.Run()
	return evalPrompt(value, prompt.Label.(string), err)
}

func evalPrompt(value, label string, err error) string {
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			print("Aborting!")
			os.Exit(1)
		} else {
			exitWithMessage(fmt.Sprintf("Specified %s is not valid.", label))
		}
	}
	return value
}

var regexMatchError = errors.New("string does not match regex")

func regexValidator(regex, str string) error {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		return err
	}
	if !match {
		return regexMatchError
	}
	return nil
}

func projectNameValidator(projectName string) error {
	err := regexValidator(project.ProjectNameRegex, projectName)
	if err != nil {
		if errors.Is(err, regexMatchError) {
			return fmt.Errorf(
				"%s is not a valid project name. It must match the regex %s",
				projectName,
				project.ProjectNameRegex,
			)
		} else {
			return fmt.Errorf("Failed to validate project name %s", projectName)
		}
	}
	return nil
}

func serviceNameValidator(serviceName string) error {
	err := regexValidator(project.ServiceNameRegex, serviceName)
	if err != nil {
		if errors.Is(err, regexMatchError) {
			return fmt.Errorf(
				"%s is not a valid service name. It must match the regex %s",
				serviceName,
				project.ServiceNameRegex,
			)
		} else {
			return fmt.Errorf("Failed to validate project name %s", serviceName)
		}
	}
	return nil
}

func updateTemplates() {
	// Get or update the templates.
	logger.Infof("Updating the templates repository...")
	out := newOut("git: ")
	err := git.UpdateTemplates(out, foldTemplates)
	exitIfErr(err, `Failed to update the template repository.
Please ensure you are connected to the internet and that you are able to access github.com`)
}
