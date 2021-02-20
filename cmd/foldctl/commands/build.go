package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build [service]",
	Example: "foldctl build ./service/",
	Short:   "Builds the specified service",
	Long: `Build the service located at the path you provide.
In order for a directory to be a valid service, it just needs a Dockerfile that implements
the fold runtime interface.

Currently, the only way to do this is by extending one of the based foldrt images. Check
out the examples or use one of the templates from the 'new' command to see how to do it.

It's worth noting that the images built with this command are not 'production ready'. They are
built for the purpose of local development and just represent the current state of the code. They
don't have a version or anything like that and this is why they aren't appropriate for production
use. The images will be tagged with the pattern foldlocal/<abs-path-hash>/<directory-name>. This ensures
that the images tags are unique in case you end up with two services with the same name for different
projects.

When you go on to deploy a service, the image will be given a tag that is tied to the specific
version of the code that built it, ensuring you can roll services back, forward, promote images
between stages, etc. Additionally, rather than using the name of the directory, the name used
will be extracted from the service manifest, and will match the service path on your gateway.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		out := newOut("docker: ")
		p := loadProjectWithRuntime(out)
		service := getService(p, path)
		service.Build(commandCtx, out)
		// TODO exit with appropriate error message
	},
}
