package runtime

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	"github.com/foldsh/fold/runtime/supervisor"
)

type RuntimeTestCase struct {
	t    *testing.T
	rt   *runtime.Runtime
	done chan struct{}
}

func NewRuntimeTestCase(t *testing.T, bin string, options ...runtime.Option) *RuntimeTestCase {
	t.Parallel()
	logger := logging.NewTestLogger()
	cmd := "go"
	args := []string{"run", bin}
	done := make(chan struct{})
	sout := &bytes.Buffer{}
	serr := &bytes.Buffer{}
	defaults := []runtime.Option{
		runtime.WithSupervisor(supervisor.NewSupervisor(logger, cmd, args, sout, serr)),
		runtime.CrashPolicy(runtime.KEEP_ALIVE),
	}
	rt := runtime.NewRuntime(
		logger,
		cmd,
		args,
		done,
		append(defaults, options...)...,
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

type ResponseWriter struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{headers: make(map[string][]string)}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.headers
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return len(b), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) Status() int {
	return rw.statusCode
}

func (rw *ResponseWriter) String() string {
	return string(rw.body)
}
