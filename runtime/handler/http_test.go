package handler_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/handler"
)

func TestServeHTTP(t *testing.T) {
	h := handler.NewHTTP(logging.NewTestLogger(), serve{}, ":12344")
	done := make(chan struct{})

	go func() {
		// Give the server a bit of time to come up
		time.Sleep(20 * time.Millisecond)
		// Assert that the server is running and that we can make a request
		resp, err := http.Get("http://localhost:12344")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if string(body) != "fold" {
			t.Fatalf("Expected fold but found %s", string(body))
		}
		// Now lets shut it all down
		h.Shutdown(context.Background(), done)
	}()

	// We serve in the main test goroutine. This will exercise the shutdown logic as the test
	// will block forever without it.
	h.Serve()
	<-done
}

type serve struct{}

func (s serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "fold")
}
