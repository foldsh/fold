package commands

import (
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Example: "foldctl version\nfoldctl --version",
		Short:   "Prints the version of foldctl",
		Long: `Prints the version of foldctl.
You can also use --version to get the same information.
	`,
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			root.SetArgs([]string{"--version"})
			root.Execute()
		},
	}
}
