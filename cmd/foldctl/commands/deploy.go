package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy [service]",
	Short: "Deploys the specified service to the fold platform.",
	Long: `Deploy the specified service.
This will build your service and then deploy it to your fold account.`,
	Run: func(cmd *cobra.Command, args []string) {
		var service string
		if len(args) > 0 {
			service = args[0]
		} else {
			service = "."
		}
		fmt.Printf("Deploying service %s", service)
	},
}
