package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	sout    = color.Output
	serr    = color.Error
	red     = color.New(color.FgRed).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()

	thisIsABug         = "This is a bug, please report it at https://github.com/foldsh/fold."
	checkPermissions   = "Please ensure that you have the relevant permissions to create files and directories there."
	servicePathInvalid = `The path you have specified is not a valid service.
Please check that it is a valid absolute or relative path to a fold service.`
	cantReachDocker = "Failed to contact the docker daemon. Please ensure it is running."
)

func print(f string, args ...interface{}) {
	fmt.Fprintf(sout, fmt.Sprintf("%s\n", f), args...)
}

func printErr(f string, args ...interface{}) {
	fmt.Fprintf(serr, fmt.Sprintf("%s\n", f), args...)
}

func exitWithMessage(lines ...string) {
	print("%s%s", red("ERROR\n\n"), red(strings.Join(lines, "\n")))
	os.Exit(1)
}

func exitWithErr(err error, lines ...string) {
	logger.Debugf("exiting with error: %v", err)
	exitWithMessage(append([]string{err.Error()}, lines...)...)
}

func exitIfErr(err error, lines ...string) {
	if err != nil {
		exitWithErr(err, lines...)
	}
}
