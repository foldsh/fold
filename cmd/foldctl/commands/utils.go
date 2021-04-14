package commands

import (
	"fmt"
	"strings"

	"github.com/foldsh/fold/ctl/output"
)

var (
	thisIsABug = output.Error(
		"this is a bug, please report it at https://github.com/foldsh/fold.",
	)

	checkPermissions = trimf(`
Please ensure that you have the relevant permissions to create files and directories there.`,
	)

	servicePathInvalid = output.Error(trimf(`
The path you have specified is not a valid service.
Please check that it is a valid absolute or relative path to a fold service.`))

	cantReachDocker = output.Error(
		"Failed to contact the docker daemon. Please ensure it is running.",
	)
)

func trimf(f string, args ...interface{}) string {
	return strings.TrimSpace(fmt.Sprintf(f, args...))
}
