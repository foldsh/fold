package output

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
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
	o.out.render(r)
	newLine(o.out)
}

// Displayf behaves the same as Display but it takes a format string and a variable number of
// arguments.
func (o *Output) Displayf(f string, args ...interface{}) {
	o.out.render(Line(fmt.Sprintf(f, args...)))
	newLine(o.out)
}

// Inform is used for secondary output of a command. The purpose of this kind of output is to
// keep the user informed about what is going on in the command. This output is not intended to
// be used as the input to other commands and so it is sent to 'err'. You might wish to use this
// to inform the user of progress, or display error messages.
func (o *Output) Inform(r Renderer) {
	o.err.render(r)
	newLine(o.err)
}

// Informf behaves the same as Inform but it takes a format string and a variable number of
// arguments.
func (o *Output) Informf(f string, args ...interface{}) {
	o.err.render(Line(fmt.Sprintf(f, args...)))
	newLine(o.err)
}

// Obtain an io.Writer that writes to the out stream.
func (o *Output) DisplayWriter(options ...option) io.Writer {
	return o.out.newWriter(options...)
}

// Obtain an io.Writer that writes to the err stream.
func (o *Output) InformWriter(options ...option) io.Writer {
	return o.err.newWriter(options...)
}

func newLine(m *multiplexer) {
	m.render(Line("\n"))
}

var (
	InputCancelled  = errors.New("input cancelled")
	InvalidInput    = errors.New("invalid input")
	RegexMatchError = errors.New("input does not match regex")
)

type Validator func(string) error

func (o *Output) Prompt(label string, validator Validator) (string, error) {
	p := promptui.Prompt{Label: label, Validate: promptui.ValidateFunc(validator)}
	value, err := p.Run()
	return processInputResult(value, err)
}

func (o *Output) Select(label string, items []string) (string, error) {
	s := promptui.Select{Label: label, Items: items}
	_, value, err := s.Run()
	return processInputResult(value, err)
}

func NoopValidator(_ string) error { return nil }

func RegexValidator(regex, str string) error {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		return err
	}
	if !match {
		return RegexMatchError
	}
	return nil
}

func processInputResult(value string, err error) (string, error) {
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return "", InputCancelled
		} else {
			return "", InvalidInput
		}
	}
	return value, nil
}
