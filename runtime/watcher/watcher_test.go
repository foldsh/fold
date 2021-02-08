package watcher_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/watcher"
)

const (
	N int = 4
)

func TestWatcher(t *testing.T) {
	testDir, err := ioutil.TempDir("", "foldWatcherTest")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(testDir)

	// Set up directory structure
	traverse(t, testDir, createDir)

	logger := logging.NewTestLogger()
	counter := &counter{}
	watcher, err := watcher.NewWatcher(logger, testDir, counter.increment)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if err := watcher.Watch(); err != nil {
		t.Fatalf("%+v", err)
	}
	traverse(t, testDir, mutateDir)
	// A small sleep to let all the events filter through. A bit rubbish
	// but we don't really need the ability to wait built into the watcher.
	// It's just for this test we need to give it a bit of time so I've
	// just hacked it like this.
	time.Sleep(100 * time.Millisecond)
	expectedNMutations := N*N + (N * N * N)
	if counter.count != expectedNMutations {
		t.Errorf("Expected %d mutations but found %d", expectedNMutations, counter.count)
	}
	watcher.Close()
}

func traverse(t *testing.T, root string, fn func(*testing.T, string)) {
	for i := 0; i < N; i++ {
		subDir := filepath.Join(root, fmt.Sprintf("%d", i))
		fn(t, subDir)
		for j := 0; j < N; j++ {
			subSubDir := filepath.Join(subDir, fmt.Sprintf("%d", j))
			fn(t, subSubDir)
		}
	}
}

func createDir(t *testing.T, path string) {
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Fatalf("%+v", err)
	}
}

func mutateDir(t *testing.T, path string) {
	for i := 0; i < N; i++ {
		file := filepath.Join(path, fmt.Sprintf("file-%d", i))
		if err := ioutil.WriteFile(file, []byte{}, 0644); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}

type counter struct {
	count int
}

func (c *counter) increment() {
	c.count++
}
