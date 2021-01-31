package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Diff(t *testing.T, expectation, actual interface{}, msg string) {
	if diff := cmp.Diff(expectation, actual); diff != "" {
		t.Errorf("%s (-want +got):\n%s", msg, diff)
	}
}
