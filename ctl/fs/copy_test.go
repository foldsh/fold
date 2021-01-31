package fs

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// This test has been adapted from the one in the docker code base.
// It basically just generates a large file structure and then copies it, then asserts
// that everything has been copied correctly and that the various file permissions etc have
// been preserved.
func TestCopyDir(t *testing.T) {
	srcDir, err := ioutil.TempDir("", "foldCopySrc")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	populateSrcDir(t, srcDir, 3)

	dstDir, err := ioutil.TempDir("", "foldCopyDst")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(srcDir)
	defer os.RemoveAll(dstDir)
	err = CopyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			t.Errorf("%+v", err)
		}
		if relPath == "." {
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)

		dstFileInfo, err := os.Lstat(dstPath)
		if err != nil {
			t.Errorf("%+v", err)
		}

		srcFileSys := f.Sys().(*syscall.Stat_t)
		dstFileSys := dstFileInfo.Sys().(*syscall.Stat_t)

		if srcFileSys.Dev == dstFileSys.Dev {
			if srcFileSys.Ino == dstFileSys.Ino {
				t.Errorf("Expected ino to be different for copied files. %s", relPath)
			}
		}

		if srcFileSys.Mode != dstFileSys.Mode {
			t.Errorf("Expected mode to be %v but found %v", srcFileSys.Mode, dstFileSys.Mode)
		}
		if srcFileSys.Uid != dstFileSys.Uid {
			t.Errorf("Expected uid to be %v but found %v", srcFileSys.Uid, dstFileSys.Uid)
		}
		if srcFileSys.Gid != dstFileSys.Gid {
			t.Errorf("Expected gid to be %v but found %v", srcFileSys.Gid, dstFileSys.Gid)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func randomMode(baseMode int) os.FileMode {
	for i := 0; i < 7; i++ {
		baseMode = baseMode | (1&rand.Intn(2))<<uint(i)
	}
	return os.FileMode(baseMode)
}

func populateSrcDir(t *testing.T, srcDir string, remainingDepth int) {
	if remainingDepth == 0 {
		return
	}
	for i := 0; i < 10; i++ {
		dirName := filepath.Join(srcDir, fmt.Sprintf("srcdir-%d", i))
		err := os.Mkdir(dirName, randomMode(0700))
		if err != nil {
			t.Fatalf("%+v", err)
		}
		populateSrcDir(t, dirName, remainingDepth-1)
	}

	for i := 0; i < 10; i++ {
		fileName := filepath.Join(srcDir, fmt.Sprintf("srcfile-%d", i))
		// Owner read bit set
		err := ioutil.WriteFile(fileName, []byte{}, randomMode(0400))
		if err != nil {
			t.Fatalf("%+v", err)
		}
	}
}
