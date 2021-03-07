package testutils

import (
	"testing"

	"github.com/foldsh/fold/manifest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Diff(t *testing.T, expectation, actual interface{}, msg string) {
	if diff := cmp.Diff(expectation, actual); diff != "" {
		t.Errorf("%s (-want +got):\n%s", msg, diff)
	}
}

func DiffManifest(t *testing.T, expectation, actual *manifest.Manifest) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(manifest.Manifest{}, manifest.Version{}, manifest.BuildInfo{}, manifest.Route{}),
	); diff != "" {
		t.Errorf("Manifest does not match exepctation(-want +got):\n%s", diff)
	}
}

func DiffFoldHTTPRequest(t *testing.T, expectation, actual *manifest.FoldHTTPRequest) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(manifest.FoldHTTPRequest{}, manifest.FoldHTTPProto{}, manifest.StringArray{}),
	); diff != "" {
		t.Errorf("FoldHTTPRequest does not match exepctation(-want +got):\n%s", diff)
	}
}

func DiffFoldHTTPResponse(t *testing.T, expectation, actual *manifest.FoldHTTPResponse) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(manifest.FoldHTTPResponse{}),
	); diff != "" {
		t.Errorf("FoldHTTPResponse does not match exepctation(-want +got):\n%s", diff)
	}
}
