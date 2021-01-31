package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the fold development server",
	Long: `Starts the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:8080.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting the fold development server")
	},
}
