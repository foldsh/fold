package commands

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
)

var (
	commandCtx context.Context
	logger     logging.Logger
	verbose    bool
	debug      bool

	rootCmd = &cobra.Command{
		Use:   "foldctl",
		Short: "Fold CLI",
		Long:  "Fold CLI",
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			// All successful commands print 'ok' in green at the end.
			print("%s", green("\nok"))
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithMessage("Failed to start foldctl.", thisIsABug)
	}
}

func init() {
	cobra.OnInitialize(func() {
		setUpLogger()
		setUpContext()
		loadConfig()
	})

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

func setUpLogger() {
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

func setUpContext() {
	commandCtx = context.Background()
	ctx, cancel := context.WithCancel(commandCtx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		logger.Debugf("listening for SIGINT")
		select {
		case <-c:
			print("Aborting!")
			cancel()
		case <-ctx.Done():
		}
	}()
}

func loadConfig() {
	err := loadConfigAtPath(foldHome)
	if err == nil {
		return
	} else if errors.Is(err, couldNotCreateDefaultConfig) {
		exitWithMessage("Failed to create the default foldctl config. Check you have permissions to write to ~/.fold/config.yaml")
	} else if errors.Is(err, couldNotReadConfigFile) {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
	} else {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
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
