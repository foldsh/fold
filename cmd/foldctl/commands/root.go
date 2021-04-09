package commands

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/version"
)

var (
	logger     logging.Logger
	commandCtx context.Context
	commandCfg *config.Config

	verbose bool
	debug   bool

	foldHome      string
	foldTemplates string

	rootCmd = &cobra.Command{
		Use:     "foldctl",
		Short:   "Fold CLI",
		Long:    "Fold CLI",
		Version: version.FoldVersion.String(),
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			// All other successful commands print 'ok' in green at the end.
			// We have one exception, which is 'version'. It's an exception so that
			// the behaviour of it is identical to --version.
			if cmd.CalledAs() == "version" {
				return
			}
			print("%s", green("\nok"))
		},
	}
)

func init() {
	setUpLogger()
	setUpContext()
	loadConfig()

	// Verbose/Debug output
	rootCmd.SetVersionTemplate(`{{printf "foldctl %s\n" .Version}}`)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&debug, "debug", "", false, "debug output - generally for debugging foldctl itself",
	)

	// Override the access token
	rootCmd.PersistentFlags().StringVarP(
		&commandCfg.AccessToken,
		"access-token",
		"t",
		commandCfg.AccessToken,
		"fold access token",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// TODO could be good to look at error types and choose behaviour
		// based on that.
		os.Exit(1)
	}
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
	// Set up foldHome path
	home, err := fs.FoldHome()
	exitIfErr(err, "Failed to locate fold home directory at ~/.fold.")
	foldHome = home
	foldTemplates = fs.FoldTemplates(foldHome)

	// Load the config from home, or create it
	cfg, err := config.Load(foldHome)
	if err == nil {
		commandCfg = cfg
		return
	} else if errors.Is(err, config.CreateConfigError) {
		exitWithMessage("Failed to create the default foldctl config. Check you have permissions to write to ~/.fold/config.yaml")
	} else if errors.Is(err, config.ReadConfigError) {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
	} else {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
	}
}
