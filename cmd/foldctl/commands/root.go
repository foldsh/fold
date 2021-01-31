package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
)

var (
	logger  logging.Logger
	verbose bool
	debug   bool
)
var rootCmd = &cobra.Command{
	Use:   "foldctl",
	Short: "Fold CLI",
	Long:  "Fold CLI",
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		// All successful commands print 'ok' in green at the end.
		print("%s", green("\nok"))
	},
}

func init() {
	cobra.OnInitialize(initialise)

	// Verbose/Debug output
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&debug, "debug", "", false, "debug output - generally for debugging foldctl itself",
	)

	// Override the access token
	rootCmd.PersistentFlags().StringP("access-token", "t", "", "fold access token")
	err := viper.BindPFlag("access-token", rootCmd.PersistentFlags().Lookup("access-token"))
	exitIfError(err, "Failed to bind the specified access token.", thisIsABug)
}

func initialise() {
	err := loadConfig()
	exitIfError(
		err,
		"Failed to load foldctl's config ~/.fold/config.yaml.",
		"Please ensure you have the relevant permissions to access files there.",
	)

	var logLevel logging.LogLevel
	if debug {
		logLevel = logging.Debug
	} else if verbose {
		logLevel = logging.Info
	} else {
		logLevel = logging.Error
	}
	l, err := logging.NewLogger(logLevel, false)
	exitIfError(err, "Failed to initialise the foldctl logger.", thisIsABug)
	logger = l
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithMessage("Failed to start foldctl.", thisIsABug)
	}
}

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

func saveProjectConfig(p *project.Project) {
	err := p.SaveConfig()
	exitIfError(
		err,
		"Failed to save fold config.",
		"Please check you have permission to write files in this directory.",
	)
}
