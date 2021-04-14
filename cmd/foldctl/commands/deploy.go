package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/output"
)

func NewDeployCmd(ctx *ctl.CmdCtx) *cobra.Command {
	return &cobra.Command{
		Use:     "deploy [service]",
		Example: "foldctl deploy ./service/",
		Short:   "Deploys the specified service to the fold platform",
		Long: `Deploy the specified service.
This will build your service and then deploy it to your fold account.`,
		Run: func(cmd *cobra.Command, args []string) {
			var service string
			if len(args) > 0 {
				service = args[0]
			} else {
				service = "."
			}
			ctx.Inform(output.Line(fmt.Sprintf("Deploying service %s", service)))
		},
	}
}
