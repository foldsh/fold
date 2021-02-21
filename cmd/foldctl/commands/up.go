package commands

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/manifest"
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
			displayServiceSummary(port, services)
		} else {
			var notAService project.NotAService
			if errors.As(err, &notAService) {
				exitWithErr(err)
			}
			exitWithMessage(thisIsABug)
		}
	},
}

// I am just doing this in here for now as it's not that clear where
// its natural home is and whether there is even any reason to have it
// represented more formally.
func displayServiceSummary(port int, services []*project.Service) {
	gatewayURL := fmt.Sprintf("http://localhost:%d", port)
	print(fmt.Sprintf("\nFold gateway is available at %s", gatewayURL))
	for _, service := range services {
		serviceURL := fmt.Sprintf("%s/%s", gatewayURL, service.Name)
		waitForHealthz(serviceURL)
		print("")
		resp, err := http.Get(fmt.Sprintf("%s/_foldadmin/manifest", serviceURL))
		exitIfErr(err)
		defer resp.Body.Close()
		m := &manifest.Manifest{}
		err = manifest.ReadJSON(resp.Body, m)
		exitIfErr(err)
		print(fmt.Sprintf("    %s is available at %s", service.Name, serviceURL))
		print(fmt.Sprintf("    %s routes:", service.Name))
		for _, route := range m.Routes {
			print(fmt.Sprintf("        %s %s%s", route.HttpMethod, serviceURL, route.PathSpec))
		}
	}
}

func waitForHealthz(serviceURL string) {
	var attempts int
	for {
		if attempts >= 10 {
			exitWithMessage("Service is not healthy, please check the container logs.")
		}
		resp, err := http.Get(fmt.Sprintf("%s/_foldadmin/healthz", serviceURL))
		exitIfErr(err)
		if resp.StatusCode == 200 {
			return
		}
		attempts += 1
		time.Sleep(100 * time.Millisecond)
	}
}
