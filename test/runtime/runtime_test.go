package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/foldsh/fold/runtime"
)

func TestBasicService(t *testing.T) {
	tc := NewRuntimeTestCase(t, "./testdata/basic/")
	q := tc.query("GET", "/hello/fold", "")
	q.expectStatus(200).expectBody(`{"greeting":"Hello, fold!"}`)
	tc.Done()
}

func TestLocalDevelopmentService(t *testing.T) {
	tc := NewRuntimeTestCase(t, "./testdata/basic/")

	q := tc.query("GET", "/hello/fold", "")
	q.expectStatus(200).expectBody(`{"greeting":"Hello, fold!"}`)

	// Now we're going to crash the process
	q = tc.query("GET", "/crash", "")
	q.expectStatus(500)

	// Not ideal but requests are not synchronised with state transitions. I.e. there is a small
	// window where a request can make it through to the old router. It would be good to sort this
	// out at some point but this state transition only affects local development so I think it's
	// ok for now.
	time.Sleep(50 * time.Millisecond)

	// After the crash we should find that the runtime is still responsive, but with the
	// default router.
	q = tc.query("GET", "/any/old/path", "")
	q.expectStatus(500).expectBody(`{"title":"Service is down"}`)

	tc.Done()
}

func TestHotReload(t *testing.T) {
	testDir := "./testdata/hot_reload"
	testFile := filepath.Join(testDir, "main.go")
	defer os.Remove(testFile)
	writeService(t, testFile, "hello")

	tc := NewRuntimeTestCase(t, testFile, runtime.WatchDir(5*time.Millisecond, testDir))

	q := tc.query("GET", "/greeting", "")
	q.expectStatus(200).expectBody(`{"msg":"hello"}`)

	writeService(t, testFile, "goodbye")

	// See the test case above for the rationale behind this sleep. It's not great but putting the
	// time in to make this unnecessary just isn't worth it right now. We would essentially need
	// to buffer incoming requests and let them flow through the system only when a router was up
	// and available. We'd need to run multiple threads inside the runtime which picked requests
	// off the queue and communicate results back to the calling goroutine via channels.
	// 1 second isn't strictly necessary but giving it some leeway makes it very reliable
	time.Sleep(1 * time.Second)

	q = tc.query("GET", "/greeting", "")
	q.expectStatus(200).expectBody(`{"msg":"goodbye"}`)
}

func writeService(t *testing.T, path, msg string) {
	code := fmt.Sprintf(`package main

import "github.com/foldsh/fold/sdks/go/fold"

func main() {
	svc := fold.NewService()
	svc.Get("/greeting", func(req *fold.Request, res *fold.Response) {
		res.StatusCode = 200
		res.Body = map[string]interface{}{"msg": "%s"}
	})
	svc.Start()
}`, msg)
	err := ioutil.WriteFile(path, []byte(code), 0644)
	if err != nil {
		t.Fatalf("Failed to write to file.")
	}
}
