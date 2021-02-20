package commands

import (
	"errors"

	"github.com/foldsh/fold/ctl/project"
	"github.com/spf13/cobra"
)

var port int

func init() {
	upCmd.PersistentFlags().IntVarP(&port, "port", "p", 6123, "development server port")
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:     "up [service]",
	Example: "foldctl up\nfoldctl up ./service-one/\nfoldctl up ./service-one/ ./service-two/",
	Short:   "Start the fold development server",
	Long: `Starts the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:6123.`,
	Run: func(cmd *cobra.Command, args []string) {
		// The current behaviour is that if no services are passed, we just start the network.
		out := newOut("docker: ")
		proj := loadProjectWithRuntime(out)
		proj.ConfigureGatewayPort(port)
		if services, err := proj.GetServices(args...); err == nil {
			if err := proj.Up(commandCtx, out, services...); err != nil {
				exitWithErr(err)
			}
		} else {
			var notAService project.NotAService
			if errors.As(err, &notAService) {
				exitWithErr(err)
			}
			exitWithMessage(thisIsABug)
		}
	},
}
