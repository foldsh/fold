package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/ctl/git"
	"github.com/foldsh/fold/ctl/project"
)

var (
	servicePath string
)

func init() {
	newCmd.Flags().StringVarP(&servicePath, "path", "p", ".", "Path to the service.")
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new [template] [language] [name]",
	Short: "Create a new fold service",
	Long: `Creates a new fold service from a template.

You must specify the template you wish to use, the language you wish to create
it for and the name you wish to give the service.

By default, the service will be created relative to the current directory. For example,

foldctl new basic js hello-world

will create a directory called hello-world that contains the basic template 
for a service written in javascript. 

Setting the --path flag will create the service relative to the specified path.
- If the directory does not exist, it will be created.
- If the directory exists, and is empty, then the command will populate that directory with the 
  specified template. 
- If the directory exists and is not empty, then the command will fail and inform you of the error.

The templates are all availble at https://github.com/foldsh/templates.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO it might be nice to move this to the prompt rather than having so many args.
		// TODO we probably want to remove the path flag actually and just always use '.' after validating we're in a fold project root.
		var (
			template = args[0]
			language = args[1]
			name     = args[2]
		)
		p := loadProject()
		// Build the absolute path to the new service.
		servicePath = filepath.Join(servicePath, name)
		absPath, err := filepath.Abs(servicePath)
		exitIfErr(err, servicePathInvalid)

		// Validate the service name
		service := project.Service{
			Name: name,
			Path: servicePath,
		}
		if err := service.Validate(); err != nil {
			exitWithErr(err)
		}

		// Update templates repoistory and validate template
		updateTemplates()
		selectedTemplate := validateTemplate(template, language)

		// Check if the directory is empty
		empty, err := fs.IsEmpty(absPath)
		if err == nil && !empty {
			exitWithMessage(
				fmt.Sprintf("The target directory %s already exists and is not empty.", absPath),
				"Please either choose a different name for your service or remove the existing directory.",
			)
		}

		// Create the path to the new service.
		logger.Debugf("Creating service directory")
		err = os.MkdirAll(absPath, DIR_PERMISSIONS)
		exitIfErr(err, "Failed to create a directory at the path you specified.", checkPermissions)

		// Copy the contents of the chosen template into the new directory.
		logger.Debugf("Copying project template to service directory")
		err = fs.CopyDir(selectedTemplate, absPath)
		if err != nil {
			os.RemoveAll(absPath)
			exitWithMessage(
				"Failed to create a project template at the path you specified.",
				checkPermissions,
			)
		}

		// Update the project config
		logger.Debugf("Adding service to project")
		p.AddService(service)
		logger.Debugf("Saving project config")
		saveProjectConfig(p)
	},
}

func updateTemplates() {
	// Get or update the templates.
	logger.Infof("Updating the templates repository...")
	out := newStreamLinePrefixer(serr, blue("git: "))
	err := git.UpdateTemplates(out, foldTemplates)
	exitIfErr(err, `Failed to update the template repository.
Please ensure you are connected to the internet and that you are able to access github.com`)
}

func validateTemplate(template, language string) string {
	// Check that the specified template is valid.
	selectedTemplate := filepath.Join(foldTemplates, template, language)
	logger.Debugf("inferred template is %s, checking if it is valid", selectedTemplate)
	stat, err := os.Stat(selectedTemplate)
	logger.Debugf("stat for inferred template path %v", stat)
	if err != nil {
		if os.IsNotExist(err) {
			// directory does not exist
			logger.Debugf("template does not exist")
			msg := fmt.Sprintf(`The specified template %s/%s does not exist.
Check the fold templates repository for available templates.
You can find a link to the repository in the help for this command, run:

foldctl new --help`, template, language)
			exitWithMessage(msg)
		} else {
			// other error
			logger.Debugf("unexpected error")
			exitWithMessage("Failed to validate the template specified.", thisIsABug)
		}
	}
	return selectedTemplate
}
