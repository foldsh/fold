package output_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/foldsh/fold/ctl/output"
)

func TestError(t *testing.T) {
	var sout bytes.Buffer
	var serr bytes.Buffer
	o := output.NewOutput(&sout, &serr)
	o.Inform(output.Error("failed"))

	assert.Equal(t, "Error: failed", serr.String(), "did not match expected output")
}
