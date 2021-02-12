package watcher_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/watcher"
)

func TestDebouncer(t *testing.T) {
	testDir, err := ioutil.TempDir("", "foldDebounce")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(testDir)

	// Set up directory structure
	traverse(t, testDir, createDir)

	logger := logging.NewTestLogger()
	counter := &counter{}
	debouncer := watcher.NewDebouncer(100*time.Millisecond, counter.increment)
	watcher, err := watcher.NewWatcher(logger, testDir, debouncer.OnChange)
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
	// Given we debounce for 100 milliseconds, we should only see counter.increment
	// actually called once, even though we know a load of events have been detected.
	expectedNMutations := 1
	if counter.count != expectedNMutations {
		t.Errorf("Expected %d mutations but found %d", expectedNMutations, counter.count)
	}
	watcher.Close()
}
