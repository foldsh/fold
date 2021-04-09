package commands

import (
	"github.com/spf13/cobra"
)

func NewDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "down [service]",
		Example: "foldctl down",
		Short:   "Stops the fold development server",
		Long: `Stops the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:6123.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			out := newOut("docker: ")
			proj := loadProjectWithRuntime(out)
			proj.Down()
			// TODO exit with appropriate error message
		},
	}
}
