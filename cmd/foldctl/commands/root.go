package commands

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		// TODO could be good to look at error types and choose behaviour
		// based on that.
		os.Exit(1)
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
	exitIfErr(err, "Failed to bind the specified access token.", thisIsABug)
}

func setUpLogger() {
	var (
		l   logging.Logger
		err error
	)
	if debug {
		l, err = logging.NewLogger(logging.Debug, false)
	} else {
		l, err = logging.NewCLILogger(logging.Info)
	}
	exitIfErr(err, "Failed to initialise the foldctl logger.", thisIsABug)
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
