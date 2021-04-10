package commands

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/version"
)

var (
	cmdctx *ctl.CmdCtx

	verbose             bool
	debug               bool
	accessTokenOverride string

	foldHome      string
	foldTemplates string

	rootCmd = &cobra.Command{
		Use:     "foldctl",
		Short:   "Fold CLI",
		Long:    "Fold CLI",
		Version: version.FoldVersion.String(),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			setUpLogger()
		},
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

func Execute() {
	// Verbose/Debug output
	rootCmd.SetVersionTemplate(`{{printf "foldctl %s\n" .Version}}`)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&debug, "debug", "", false, "debug output - generally for debugging foldctl itself",
	)

	// Override the access token
	rootCmd.PersistentFlags().StringVarP(
		&accessTokenOverride,
		"access-token",
		"t",
		"",
		"fold access token",
	)

	// Create the Context and send up a listener for SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			print("Aborting!")
			cancel()
		case <-ctx.Done():
		}
	}()

	// Set up the CLI config
	cfg := loadConfig()

	// Create the new CmdCtx
	// TODO passing the logger as nil because the debug flag is not bound before we set it up.
	// However I want to be able to pass the context as a dependency to the commands so there is
	// a bit of a chicken and egg problem... For now I'm just using a PreRun hook to set the logger
	// once the flags are bound.
	cmdctx = ctl.NewCmdCtx(ctx, nil, cfg, serr)
	addCommands(cmdctx, rootCmd)

	cobra.OnInitialize(func() {
		setUpLogger()
	})

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
	cmdctx.Logger = l
}

func loadConfig() *config.Config {
	// Set up foldHome path
	home, err := fs.FoldHome()
	exitIfErr(err, "Failed to locate fold home directory at ~/.fold.")
	foldHome = home
	foldTemplates = fs.FoldTemplates(foldHome)

	// Load the config from home, or create it
	cfg, err := config.Load(foldHome)
	if err == nil {
		if accessTokenOverride != "" {
			cfg.AccessToken = accessTokenOverride
		}
		return cfg
	} else if errors.Is(err, config.CreateConfigError) {
		exitWithMessage("Failed to create the default foldctl config. Check you have permissions to write to ~/.fold/config.yaml")
	} else if errors.Is(err, config.ReadConfigError) {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
	} else {
		exitWithMessage("Failed to read the foldctl config file at ~/.fold/config.yaml. Please ensure it is valid yaml.")
	}
	return nil
}

func addCommands(ctx *ctl.CmdCtx, root *cobra.Command) {
	root.AddCommand(NewVersionCmd(ctx))
	root.AddCommand(NewBuildCmd(ctx))
	root.AddCommand(NewUpCmd(ctx))
	root.AddCommand(NewDeployCmd(ctx))
	root.AddCommand(NewDownCmd(ctx))
	root.AddCommand(NewNewCmd(ctx))
}
