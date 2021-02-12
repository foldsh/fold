package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl/project"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new fold project in the current directory.",
	Long: `Initialises a new fold project in the current directory.
The command will start a series of prompts to guide you through the set up process.`,
	Run: func(cmd *cobra.Command, args []string) {
		if project.IsAFoldProject(".") {
			exitWithMessage("Already a fold project.")
		}
		projectName := stringPrompt("Project name")
		maintainer := stringPrompt("Maintainer")
		email := stringPrompt("Email")
		repo := stringPrompt("Repository")

		p := &project.Project{
			Name:       projectName,
			Maintainer: maintainer,
			Email:      email,
			Repository: repo,
			Services:   []*project.Service{},
		}
		p.ConfigureLogger(logger)
		saveProjectConfig(p)
	},
}

func stringPrompt(label string) string {
	prompt := promptui.Prompt{Label: label}
	value, err := prompt.Run()
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
