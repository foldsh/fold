package output

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithPrefix(t *testing.T) {
	out := &testWriter{}
	m := newMultiplexer(out)
	o := m.newWriter(WithPrefix("TEST: "))

	writeBytesAndTestN(t, o, "this is line 1\nthis is line 2\rthis is")
	writeBytesAndTestN(t, o, " line 3\nthis is line 4\n")
	writeBytesAndTestN(t, o, "this is line 5")
	writeBytesAndTestN(t, o, "\n")
	writeBytesAndTestN(t, o, "\n")

	expectation := []string{
		"TEST: this is line 1\n",
		"TEST: this is line 2\r",
		"TEST: this is line 3\n",
		"TEST: this is line 4\n",
		"TEST: this is line 5\n",
		"TEST: \n",
	}
	assert.Equal(t, expectation, out.lines)
}

func TestConcurrentWrites(t *testing.T) {
	out := &testWriter{}
	m := newMultiplexer(out)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			o := m.newWriter(WithPrefix("TEST: "))
			writeBytesAndTestN(t, o, "first line\nsecond")
			writeBytesAndTestN(t, o, " line\n")
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equalf(
		t,
		200,
		len(out.lines),
		"Expected to see 200 lines written but found %d",
		len(out.lines),
	)
	for _, line := range out.lines {
		assert.Truef(
			t,
			strings.HasPrefix(line, "TEST: "),
			"All lines should start with 'TEST: ' but found %s",
			line,
		)
		assert.Truef(
			t,
			strings.HasSuffix(line, "line\n"),
			"All lines should end with 'line\\n' but found %s",
			line,
		)
	}
}

func TestWriteComplexData(t *testing.T) {
	out := &testWriter{}
	m := newMultiplexer(out)

	file, err := os.Open("./testdata/logs.txt")
	require.Nil(t, err)

	scanner := bufio.NewScanner(file)
	// We'll scan each byte at a time and write them all to simulate a stream a bit better.
	scanner.Split(bufio.ScanBytes)
	var bytes []byte
	w := m.newWriter()
	for scanner.Scan() {
		b := scanner.Bytes()
		bytes = append(bytes, b[0])
		w.Write(b)
	}
	assert.Equal(t, string(bytes), strings.Join(out.lines, ""))
}

type testWriter struct {
	lines []string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.lines = append(tw.lines, string(p))
	return
}

func writeBytesAndTestN(t *testing.T, w io.Writer, input string) {
	bytes := []byte(input)
	n, _ := w.Write(bytes)
	assert.Equal(t, len(bytes), n)
}
