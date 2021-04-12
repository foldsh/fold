package output_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/foldsh/fold/ctl/output"
)

func TestError(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Inform(output.Error("failed"))

	assert.Equal(t, "\nError: failed\n", serr.String(), "did not match expected output")
}

func TestSuccess(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Inform(output.Success("success"))

	assert.Equal(t, "\nsuccess\n", serr.String(), "did not match expected output")
}

type testRenderer struct {
	a int
	b int
}

func (t testRenderer) Render() string {
	return fmt.Sprintf("%d%d", t.a, t.b)
}

func TestStructRenderer(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Inform(testRenderer{1, 2})

	assert.Equal(t, "12\n", serr.String(), "did not match expected output")
}

func TestDisplayf(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Displayf("foo %s %s %d", "bar", "baz", 42)

	assert.Equal(t, "foo bar baz 42\n", sout.String(), "did not match expected output")
}

func TestInformf(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Informf("foo %s %s %d", "bar", "baz", 42)

	assert.Equal(t, "foo bar baz 42\n", serr.String(), "did not match expected output")
}

func TestWriters(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)

	outw := o.DisplayWriter(output.WithPrefix("out: "))
	errw := o.InformWriter(output.WithPrefix("err: "))

	outw.Write([]byte("display\n"))
	errw.Write([]byte("inform\n"))

	assert.Equal(t, "out: display\n", sout.String(), "did not match expected output")
	assert.Equal(t, "err: inform\n", serr.String(), "did not match expected output")
}
