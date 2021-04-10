package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/foldsh/fold/cmd/foldctl/commands"
	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/logging"
)

func main() {
	// Set up the logger
	logger := setUpLogger()

	// Create the Context and listen for an interrupt
	ctx, cleanup := createContext()
	defer cleanup()

	// Load the CLI config
	cfg := loadConfig()

	// Create the new CmdCtx
	cmdctx := ctl.NewCmdCtx(ctx, logger, cfg, os.Stderr)

	// Create the root command
	cmd := commands.NewRootCmd(cmdctx)

	// And execute it
	if err := cmd.Execute(); err != nil {
		// TODO could be good to look at error types and choose behaviour
		// based on that.
		os.Exit(1)
	}
}

func createContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
			return
		}
	}()
	return ctx, func() {
		signal.Stop(sigChan)
		cancel()
	}
}

func setUpLogger() logging.Logger {
	var (
		logger logging.Logger
		err    error
	)
	debug := os.Getenv("DEBUG")
	if debug == "1" || strings.ToLower(debug) == "true" {
		logger, err = logging.NewLogger(logging.Debug, false)
	} else {
		logger, err = logging.NewCLILogger(logging.Info)
	}
	if err != nil {
		fmt.Println("Failed to initialise the logger")
		os.Exit(1)
	}
	return logger
}

func loadConfig() *config.Config {
	// Set up foldHome path
	home, err := fs.FoldHome()
	if err != nil {
		fmt.Println("Failed to locate fold home directory at ~/.fold.")
		os.Exit(1)
	}

	// Load the config from home, or create it
	cfg, err := config.Load(home)
	if err == nil {
		return cfg
	} else {
		var msg string
		if errors.Is(err, config.CreateConfigError) {
			msg = "Failed to create the default foldctl config.\nCheck you have permissions to write to ~/.fold/config.yaml"
		} else {
			msg = "Failed to read the foldctl config file at ~/.fold/config.yaml.\nPlease ensure it is valid yaml."
		}
		fmt.Println(msg)
		os.Exit(1)
	}
	return nil
}
