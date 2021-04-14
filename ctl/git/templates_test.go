package git_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/foldsh/fold/ctl/git"
)

func TestCloneTemplates(t *testing.T) {
	dir, err := ioutil.TempDir("", "fold.clone-templates.test")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(dir)
	templatesDir := filepath.Join(dir, "templates")

	out := &bytes.Buffer{}
	// The first time it should clone
	err = git.UpdateTemplates(out, templatesDir, "v0.1.2")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if _, err := os.Stat(filepath.Join(templatesDir, "README.md")); err != nil {
		t.Fatalf("%+v", err)
	}

	// The second time it should just update without error.
	err = git.UpdateTemplates(out, templatesDir, "v0.1.2")
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
