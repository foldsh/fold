package runtime

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	"github.com/foldsh/fold/runtime/supervisor"
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

type RuntimeTestCase struct {
	t    *testing.T
	rt   *runtime.Runtime
	done chan struct{}
}

func NewRuntimeTestCase(t *testing.T, bin string) *RuntimeTestCase {
	t.Parallel()
	logger := logging.NewTestLogger()
	cmd := "go"
	args := []string{"run", bin}
	done := make(chan struct{})
	sout := &bytes.Buffer{}
	serr := &bytes.Buffer{}
	rt := runtime.NewRuntime(
		logger,
		cmd,
		args,
		done,
		runtime.WithSupervisor(supervisor.NewSupervisor(logger, cmd, args, sout, serr)),
		runtime.CrashPolicy(runtime.KEEP_ALIVE),
	)
	rt.Start()
	return &RuntimeTestCase{
		t:    t,
		rt:   rt,
		done: done,
	}
}

func (r *RuntimeTestCase) Done() {
	r.rt.Stop()
	<-r.done
}

func (r *RuntimeTestCase) query(method, path, body string) *QueryAssertion {
	w := NewResponseWriter()
	req, err := http.NewRequest(
		method,
		path,
		ioutil.NopCloser(strings.NewReader(body)),
	)
	if err != nil {
		r.t.Errorf("%+v", err)
		r.rt.Stop()
		<-r.done
	}
	r.rt.ServeHTTP(w, req)
	return &QueryAssertion{r.t, w}
}

type QueryAssertion struct {
	t *testing.T
	w *ResponseWriter
}

func (qa *QueryAssertion) expectStatus(status int) *QueryAssertion {
	if qa.w.Status() != status {
		qa.t.Errorf("Expected a %d response code but found %d", status, qa.w.Status())
	}
	return qa
}

func (qa *QueryAssertion) expectBody(body string) *QueryAssertion {
	actual := qa.w.String()
	if actual != body {
		testutils.Diff(qa.t, body, actual, "Body did not match expectation")
	}
	return qa
}
