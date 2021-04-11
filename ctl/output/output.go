package output

import (
	"io"

	"github.com/fatih/color"
)

type Output struct {
	out *multiplexer
	err *multiplexer
}

func NewOutput(out io.Writer, err io.Writer) *Output {
	return &Output{out: newMultiplexer(out), err: newMultiplexer(err)}
}

func NewColorOutput() *Output {
	return NewOutput(color.Output, color.Error)
}

// Display is used for the primary output of a command. It is always sent to 'out' and can
// be used to easily pipe output into other commands, files, etc. You might wish to use this
// to display a table that describes resources, for example.
func (o *Output) Display(r Renderer) {
}

// Displayf behaves the same as Display but it takes a format string and a variable number of
// arguments.
func (o *Output) Displayf(f string, args ...interface{}) {
}

// Inform is used for secondary output of a command. The purpose of this kind of output is to
// keep the user informed about what is going on in the command. This output is not intended to
// be used as the input to other commands and so it is sent to 'err'. You might wish to use this
// to inform the user of progress, or display error messages.
func (o *Output) Inform(r Renderer) {
}

// Informf behaves the same as Inform but it takes a format string and a variable number of
// arguments.
func (o *Output) Informf(f string, args ...interface{}) {
}
