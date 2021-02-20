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
	"github.com/foldsh/fold/ctl/project"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialise a new fold project.",
	Long: fmt.Sprintf(
		`Initialises a new fold project in the current directory or in the specified path.

If you do not specify a path then the project is initialised in the current directory and you will 
be prompted for a project name.

If you specify a path then the end of the path will be selected as the project name. 
For example, running:

foldctl init fold/my-project

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
			projectName = stringPrompt(promptui.Prompt{Label: prompt, Validate: projectNameValidator})
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
		maintainer := stringPrompt(promptui.Prompt{Label: "Maintainer"})
		email := stringPrompt(promptui.Prompt{Label: "Email"})
		repo := stringPrompt(promptui.Prompt{Label: "Repository"})

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

func stringPrompt(prompt promptui.Prompt) string {
	value, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			print("Aborting!")
			os.Exit(1)
		} else {
			exitWithMessage(fmt.Sprintf("Specified %s is not valid.", prompt.Label))
		}
	}
	return value
}

func projectNameValidator(projectName string) error {
	match, err := regexp.MatchString(project.ProjectNameRegex, projectName)
	if err != nil {
		return fmt.Errorf("Failed to validate project name %s", projectName)
	}
	if !match {
		return fmt.Errorf(
			"%s is not a valid project name. It must match the regex %s",
			projectName,
			project.ProjectNameRegex,
		)
	}
	return nil
}
