package output_test

import (
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/internal/testutils"
)

func TestWithPrefix(t *testing.T) {
	out := &testWriter{}
	m := output.NewMultiplexer(out)
	o := m.Output(output.WithPrefix("TEST: "))

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
	testutils.Diff(t, expectation, out.lines, "Written lines did not match expectation")
}

func TestConcurrentWrites(t *testing.T) {
	out := &testWriter{}
	m := output.NewMultiplexer(out)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			o := m.Output(output.WithPrefix("TEST: "))
			writeBytesAndTestN(t, o, "first line\nsecond")
			writeBytesAndTestN(t, o, " line\n")
			wg.Done()
		}()
	}
	wg.Wait()
	if len(out.lines) != 200 {
		t.Errorf("Expected to see 200 lines written but found %d", len(out.lines))
	}
	for _, line := range out.lines {
		if !strings.HasPrefix(line, "TEST: ") {
			t.Errorf("All lines should start with 'TEST: ' but found %s", line)
		}
		if !strings.HasSuffix(line, "line\n") {
			t.Errorf("All lines should end with 'line\\n' but found %s", line)
		}
	}
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
	if n != len(bytes) {
		t.Errorf("Expected to write %d bytes but wrote %d", len(bytes), n)
	}
}
