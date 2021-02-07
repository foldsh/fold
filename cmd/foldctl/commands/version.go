package commands

import (
	"fmt"

	"github.com/foldsh/fold/ctl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of foldctl.",
	Long: `Prints the version of foldctl.
You can also use -V, or --version to get the same information.
	`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ctl.FoldctlVersion.String())
	},
}
