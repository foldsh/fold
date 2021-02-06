package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up [service]",
	Short: "Start the fold development server",
	Long: `Starts the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:8080.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// The current behaviour is that if no services are passed, we just start the network.
		// TODO accept any number of services and pass the slice onto proj.Up
		print("Starting the fold development server...")
		// var servicePath string
		// if len(args) == 1 {
		// 	servicePath = args[0]
		// }
		out := newOut("docker: ")
		proj := loadProjectWithRuntime(out)
		services := proj.GetServices(args...)
		proj.Up(commandCtx, out, services...)
		// TODO exit with appropriate error message
	},
}
