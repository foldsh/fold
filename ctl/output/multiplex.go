package output

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// This lets you write to an io.Writer with the usual io.Writer interface. However it will also
// detect new lines in the output and add a specified prefix to each line.
type multiplexer struct {
	output io.Writer
	mu     *sync.Mutex
}

func newMultiplexer(output io.Writer) *multiplexer {
	return &multiplexer{output: output, mu: &sync.Mutex{}}
}

func (m *multiplexer) Display(r Renderer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Fprint(m.output, r.Render())
}

func (m *multiplexer) Output(options ...option) io.Writer {
	out := &output{m: m}
	for _, o := range options {
		o(out)
	}
	return out
}

type output struct {
	m   *multiplexer
	buf bytes.Buffer

	// Set by options
	prefix string
}

// The procedure here is fairly simple, but a little fiddly.
// Basically, we read the input a byte at a time, buffering it all.
// When we encounter a new line, we flush the buffered input to the specified output.
// We then carry on reading until we run out of input.
func (out *output) Write(p []byte) (n int, err error) {
	for _, b := range p {
		// Always write the byte as we want to write the new lines too.
		// We can ignore the error here, the docs state it is always nil, this panics if something
		// goes wrong https://golang.org/pkg/bytes/#Buffer.WriteByte.
		out.buf.WriteByte(b)
		n += 1
		if b == '\n' || b == '\r' {
			out.m.Display(Line(fmt.Sprintf("%s%s", out.prefix, out.buf.String())))
			out.buf.Reset()
			continue
		}
	}
	return n, err
}

type option func(*output)

func WithPrefix(prefix string) option {
	return func(out *output) {
		out.prefix = prefix
	}
}
