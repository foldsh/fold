package commands

import (
	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/output"
)

func NewDownCmd(ctx *ctl.CmdCtx) *cobra.Command {
	return &cobra.Command{
		Use:     "down [service]",
		Example: "foldctl down",
		Short:   "Stops the fold development server",
		Long: `Stops the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:6123.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			out := ctx.InformWriter(output.WithPrefix(output.Blue("docker: ")))
			proj := loadProjectWithRuntime(ctx, out)
			proj.Down()
		},
	}
}
