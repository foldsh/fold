package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops the fold development server",
	Long: `Stops the fold development server.
This will build all of your services and wire them up to a local gateway you can
access on http://localhost:8080.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting the fold development server")
	},
}
