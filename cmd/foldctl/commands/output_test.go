package commands

import (
	"testing"

	"github.com/foldsh/fold/internal/testutils"
)

func TestStreamLinePrefixer(t *testing.T) {
	out := &testWriter{}
	slp := newStreamLinePrefixer(out, "TEST: ")

	writeBytesAndTestN(t, slp, "this is line 1\nthis is line 2\rthis is")
	writeBytesAndTestN(t, slp, " line 3\nthis is line 4\n")
	writeBytesAndTestN(t, slp, "this is line 5")
	writeBytesAndTestN(t, slp, "\n")
	writeBytesAndTestN(t, slp, "\n")

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

type testWriter struct {
	lines []string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.lines = append(tw.lines, string(p))
	return
}

func writeBytesAndTestN(t *testing.T, slp *streamLinePrefixer, input string) {
	bytes := []byte(input)
	n, _ := slp.Write(bytes)
	if n != len(bytes) {
		t.Errorf("Expected to write %d bytes but wrote %d", len(bytes), n)
	}
}
