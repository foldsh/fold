package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down [service]",
	Short: "Stops the fold development server",
	Long: `Stops the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:8080.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		print("Stopping the fold development server...")

		out := newOut("docker: ")
		proj := loadProjectWithRuntime(out)
		proj.Down()
		// TODO exit with appropriate error message
	},
}
