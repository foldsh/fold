package commands

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/manifest"
)

var (
	// Flags
	port   int
	detach bool

	// Help text
	exampleText = trimf(`
# Start the gateway
foldctl up

# Start a single service
foldctl up ./service-one/

# Start multiple services
foldctl up ./service-one/ ./service-two/"
`)

	longText = trimf(`
Starts the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:6123.
`)
)

func NewUpCmd(ctx *ctl.CmdCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up [service]",
		Short:   "Start the fold development server",
		Long:    longText,
		Example: exampleText,
		Run: func(cmd *cobra.Command, args []string) {
			// The current behaviour is that if no services are passed, we just start the network.
			out := ctx.InformWriter(output.WithPrefix(output.Blue("docker: ")))
			proj := loadProjectWithRuntime(ctx, out)
			proj.ConfigureGatewayPort(port)

			if services, err := proj.GetServices(args...); err == nil {
				if err := proj.Up(out, services...); err != nil {
					ctx.Inform(output.Error(err.Error()))
					os.Exit(1)
				}
				displayServiceSummary(ctx, port, services)
				if !detach {
					runInForeground(ctx, services)
				}
			} else {
				var notAService project.NotAService
				if errors.As(err, &notAService) {
					ctx.Inform(output.Error(err.Error()))
					os.Exit(1)
				}
				ctx.Inform(thisIsABug)
				os.Exit(1)
			}
		},
	}
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 6123, "development server port")
	cmd.PersistentFlags().BoolVarP(&detach, "detach", "d", false, "run in the background")
	return cmd
}

func runInForeground(
	ctx *ctl.CmdCtx,
	services []*project.Service,
) {
	var streams []*container.LogStream

	for _, service := range services {
		go func() {
			ctx.Logger.Debugf("Listening to logs for service %s", service.Name)
			prefix := output.Blue(fmt.Sprintf("%s: ", service.Name))
			out := ctx.DisplayWriter(output.WithPrefix(prefix))
			ls, err := service.Logs()
			if err != nil {
				ctx.Inform(output.Error(err.Error()))
				return
			}
			streams = append(streams, ls)
			if err := ls.Stream(out); err != nil {
				ctx.Inform(output.Error(err.Error()))
				return
			}
		}()
	}
	<-ctx.Done()
	ctx.Logger.Debugf("SIGINT received by context")
	for _, stream := range streams {
		stream.Stop()
	}
}

// I am just doing this in here for now as it's not that clear where
// its natural home is and whether there is even any reason to have it
// represented more formally.
func displayServiceSummary(ctx *ctl.CmdCtx, port int, services []*project.Service) {
	gatewayURL := fmt.Sprintf("http://localhost:%d", port)
	ctx.Informf("\nFold gateway is available at %s", gatewayURL)
	for _, service := range services {
		serviceURL := fmt.Sprintf("%s/%s", gatewayURL, service.Name)
		waitForHealthz(ctx, serviceURL)
		ctx.Informf("")
		resp, err := http.Get(fmt.Sprintf("%s/_foldadmin/manifest", serviceURL))
		if err != nil {
			ctx.Inform(output.Error(err.Error()))
		}
		defer resp.Body.Close()
		m := &manifest.Manifest{}
		err = manifest.ReadJSON(resp.Body, m)
		if err != nil {
			ctx.Inform(output.Error(err.Error()))
		}
		ctx.Informf("    %s is available at %s", service.Name, serviceURL)
		ctx.Informf("    %s routes:", service.Name)
		for _, route := range m.Routes {
			ctx.Informf("        %s %s%s", route.HttpMethod, serviceURL, route.Route)
		}
	}
}

func waitForHealthz(ctx *ctl.CmdCtx, serviceURL string) {
	var attempts int
	healthz := fmt.Sprintf("%s/_foldadmin/healthz", serviceURL)
	for {
		if attempts >= 10 {
			ctx.Inform(output.Error("service is not healthy"))
			ctx.Inform(output.Line("Please check the container logs."))
			os.Exit(1)
		}
		resp, err := http.Get(healthz)
		if resp != nil && err == nil {
			if resp.StatusCode == 200 {
				return
			}
		}
		attempts += 1
		time.Sleep(100 * time.Millisecond)
	}
}
