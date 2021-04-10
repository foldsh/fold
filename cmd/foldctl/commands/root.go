package commands

import (
	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/version"
)

var (
	verbose             bool
	accessTokenOverride string
)

func NewRootCmd(ctx *ctl.CmdCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "foldctl",
		Short:   "Fold CLI",
		Long:    "Fold CLI",
		Version: version.FoldVersion.String(),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			if accessTokenOverride != "" {
				ctx.Config.AccessToken = accessTokenOverride
			}
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
	// Verbose/Debug output
	cmd.SetVersionTemplate(`{{printf "foldctl %s\n" .Version}}`)
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Override the access token
	cmd.PersistentFlags().StringVarP(
		&accessTokenOverride,
		"access-token",
		"t",
		"",
		"fold access token",
	)

	cmd.AddCommand(NewVersionCmd(ctx))
	cmd.AddCommand(NewBuildCmd(ctx))
	cmd.AddCommand(NewUpCmd(ctx))
	cmd.AddCommand(NewDeployCmd(ctx))
	cmd.AddCommand(NewDownCmd(ctx))
	cmd.AddCommand(NewNewCmd(ctx))

	return cmd
}
