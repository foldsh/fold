package commands

import (
	"bytes"
	"fmt"
	"io"
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
)

func print(f string, args ...interface{}) {
	fmt.Fprintf(sout, fmt.Sprintf("%s\n", f), args...)
}

func printErr(f string, args ...interface{}) {
	fmt.Fprintf(serr, fmt.Sprintf("%s\n", f), args...)
}

func exitIfError(err error, lines ...string) {
	if err != nil {
		exitWithMessage(lines...)
	}
}

func exitWithMessage(lines ...string) {
	print("%s%s", red("Error\n\n"), red(strings.Join(lines, "\n")))
	os.Exit(1)
}

// This lets you write to an io.Writer with the usual io.Writer interface. However it will also
// detect new lines in the output and add a specified prefix to each line.
type streamLinePrefixer struct {
	output io.Writer
	prefix string
	buf    bytes.Buffer
}

func newStreamLinePrefixer(output io.Writer, prefix string) *streamLinePrefixer {
	return &streamLinePrefixer{output: output, prefix: prefix}
}

// The procedure here is fairly simple, but a little fiddly.
// Basically, we read the input a byte at a time, buffering it all.
// When we encounter a new line, we flush the buffered input to the specified output.
// We then carry on reading until we run out of input.
func (slp *streamLinePrefixer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		// Always write the byte as we want to write the new lines too.
		// We can ignore the error here, the docs state it is always nil, this panics if something
		// goes wrong https://golang.org/pkg/bytes/#Buffer.WriteByte.
		slp.buf.WriteByte(b)
		n += 1
		if b == '\n' || b == '\r' {
			// flush
			fmt.Fprintf(slp.output, "%s%s", slp.prefix, slp.buf.String())
			slp.buf.Reset()
			continue
		}
	}
	return n, err
}
