package commands

import (
	"fmt"

	"github.com/foldsh/fold/ctl/container"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [service]",
	Short: "Builds the specified service.",
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
		p := loadProject()
		service, err := p.GetService(path)
		exitIfError(
			err,
			fmt.Sprintf("The path %s is not a registered service.", path),
			"Please check the path you typed or, if this is a mistake, make sure that the service",
			"is registered in your fold.yaml file.",
		)
		// absPath, err := filepath.Abs(service.Path)
		absPath, err := service.AbsPath()
		exitIfError(err, servicePathInvalid)
		logger.Debugf("absolute path to service inferred as %s", absPath)
		tag := fmt.Sprintf("foldlocal/%s/%s", service.Id(), service.Name)
		print("Preparing to build service %s with tag %s", service.Name, tag)
		buildSpec := &container.BuildSpec{
			Src:    absPath,
			Image:  tag,
			Logger: logger,
			Out:    newStreamLinePrefixer(serr, blue("docker: ")),
		}
		err = container.Build(commandCtx, buildSpec)
		exitIfError(
			err,
			"Failed to build the service.",
			"Check the build logs above for more information on why this happened.",
		)
	},
}
