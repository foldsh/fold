package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/manifest"
	"github.com/spf13/cobra"
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
			out := ctx.Output(output.WithPrefix(blue("docker: ")))
			proj := loadProjectWithRuntime(ctx, out)
			proj.ConfigureGatewayPort(port)

			if services, err := proj.GetServices(args...); err == nil {
				if err := proj.Up(ctx.Context, out, services...); err != nil {
					exitWithErr(err)
				}
				displayServiceSummary(port, services)
				if !detach {
					runInForeground(ctx, services)
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
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 6123, "development server port")
	cmd.PersistentFlags().BoolVarP(&detach, "detach", "d", false, "run in the background")
	return cmd
}

func runInForeground(
	ctx *ctl.CmdCtx,
	services []*project.Service,
) {
	var chans []chan struct{}
	for _, service := range services {
		c := make(chan struct{})
		chans = append(chans, c)
		go func() {
			ctx.Debugf("Listening to logs for service %s", service.Name)
			out := ctx.Output(output.WithPrefix(fmt.Sprintf("%s: ", service.Name)))
			rc, err := service.Logs()
			if err != nil {
				// TODO output as error
				out.Write([]byte(err.Error()))
				return
			}
			buf := make([]byte, 1024)
			for {
				select {
				case <-c:
					ctx.Debugf("SIGINT received by goroutine")
					break
				default:
					n, err := rc.Read(buf)
					if err == io.EOF {
						break
					}
					out.Write(buf[:n])
				}
			}
		}()
	}
	<-ctx.Done()
	ctx.Debugf("SIGINT received by context")
	for _, c := range chans {
		close(c)
	}
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
		print("get manifest res: %v", resp)
		exitIfErr(err)
		defer resp.Body.Close()
		m := &manifest.Manifest{}
		err = manifest.ReadJSON(resp.Body, m)
		exitIfErr(err)
		print(fmt.Sprintf("    %s is available at %s", service.Name, serviceURL))
		print(fmt.Sprintf("    %s routes:", service.Name))
		for _, route := range m.Routes {
			print(fmt.Sprintf("        %s %s%s", route.HttpMethod, serviceURL, route.Route))
		}
	}
}

func waitForHealthz(serviceURL string) {
	var attempts int
	healthz := fmt.Sprintf("%s/_foldadmin/healthz", serviceURL)
	for {
		if attempts >= 10 {
			exitWithMessage("Service is not healthy, please check the container logs.")
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
