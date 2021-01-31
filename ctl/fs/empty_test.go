package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	dir, err := ioutil.TempDir("", "testIsEmpty")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(dir)
	empty, err := IsEmpty(dir)
	if empty != true {
		t.Errorf("Expected newly created tempdir to be empty but it was not.")
	}

	fileName := filepath.Join(dir, "afile.txt")
	err = ioutil.WriteFile(fileName, []byte("Hello, World!\n"), 0644)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	empty, err = IsEmpty(dir)
	if empty != false {
		t.Errorf("After creating a file in tempdir expected it to not be empty but it was.")
	}
}

func TestIsEmptyReturnsNotExistErr(t *testing.T) {
	empty, err := IsEmpty("/this/is/not/a/real/directory")
	if err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("IsEmpty should have failed with an ErrNotExist but got: %+v", err)
		}
	}
	if empty != false {
		t.Errorf("Empty should be false for directories that do not exist.")
	}
}
