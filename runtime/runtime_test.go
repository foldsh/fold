package runtime_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/router"
	"github.com/foldsh/fold/runtime/supervisor"
)

func TestStartFromDOWNState(t *testing.T) {
	// Starting the runtime should start the supervisor and then the client.
	// The router should then configure itself by calling GetManifest
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	if ctx.runtime.State() != runtime.UP {
		t.Errorf(
			"After succesfully starting the runtime should be in the UP state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestStartFromUPState(t *testing.T) {
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	// Ok, the runtime is UP so lets issue a second 'START'. This is expected to have the effect
	// of a restart.
	ctx.expectRuntimeStopTrace()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	if ctx.runtime.State() != runtime.UP {
		t.Errorf(
			"After succesfully starting the runtime should be in the UP state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestOnProcessEndCallback(t *testing.T) {
	var ended bool
	ctx := makeRuntime(t, runtime.OnProcessEnd(func() { ended = true }))
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	// The callback is asynchronous so we need to sleep to let it run.
	time.Sleep(10 * time.Millisecond)
	if ended != true {
		t.Errorf("Expected OnProcessEnd callback to be called but it wasn't")
	}
}

func TestStopFromDOWNState(t *testing.T) {
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.runtime.Stop()

	// An EXIT should close the channel
	<-ctx.done

	if ctx.runtime.State() != runtime.EXITED {
		t.Errorf(
			"After stopping the runtime should be in the EXITED state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestStopFromUPState(t *testing.T) {
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.expectRuntimeStopTrace()
	ctx.runtime.Start()
	ctx.runtime.Stop()

	// An EXIT should close the channel
	<-ctx.done

	if ctx.runtime.State() != runtime.EXITED {
		t.Errorf(
			"After stopping the runtime should be in the EXITED state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestExitOnCrash(t *testing.T) {
	// The default behaviour is simply to exit on a crash.
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()
	ctx.runtime.Emit(runtime.CRASH)

	if ctx.runtime.State() != runtime.EXITED {
		t.Errorf(
			"Expect the runtime to transition to the EXITED state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestKeepAliveOnCrash(t *testing.T) {
	// When we set the KEEP_ALIVE crash policy then a crash should transition us to the down
	// state instead.
	ctx := makeRuntime(t, runtime.CrashPolicy(runtime.KEEP_ALIVE))
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()
	ctx.runtime.Emit(runtime.CRASH)

	if ctx.runtime.State() != runtime.DOWN {
		t.Errorf(
			"Expected the runtime to transition to the DOWN state, but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestStopOnSignal(t *testing.T) {
	// We set the keep alive policy as it is only with that setting that this test is interesting.
	ctx := makeRuntime(t, runtime.CrashPolicy(runtime.KEEP_ALIVE))
	defer ctx.Finish()

	mockSignal := make(chan struct{})

	// This is the start
	ctx.supervisor.EXPECT().Start(map[string]string{"FOLD_SOCK_ADDR": SOCKET})
	ctx.supervisor.EXPECT().Wait().DoAndReturn(func() error {
		// Sleep for a bit to simulate the process running
		<-mockSignal
		return supervisor.TerminatedBySignal
	})
	ctx.client.EXPECT().Start(SOCKET)
	ctx.client.EXPECT().GetManifest(gomock.Any())
	ctx.router.EXPECT().Configure(gomock.Any())

	ctx.runtime.Start()

	// If we close the channel we should see a stop trace.
	ctx.expectRuntimeStopTrace()
	close(mockSignal)

	<-ctx.done

	if ctx.runtime.State() != runtime.EXITED {
		t.Errorf(
			"Expected the runtime to transition to EXIT but found %v",
			ctx.runtime.State(),
		)
	}
}

func TestDoRequestInUPState(t *testing.T) {
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	// In the UP state, we expect the request to be passed through to the actual application router.
	ctx.router.EXPECT().ServeHTTP(rw, req)
	ctx.runtime.ServeHTTP(rw, req)
}

func TestDoRequestInDOWNState(t *testing.T) {
	// The handler should result in the client getting called.
	ctx := makeRuntime(t)
	defer ctx.Finish()

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	// In the DOWN state we expect the request to be passed on to the default router.
	ctx.defaultRouter.EXPECT().ServeHTTP(rw, req)
	ctx.runtime.ServeHTTP(rw, req)
}

func TestDoRequestInEXITEDState(t *testing.T) {
	// The handler should result in the client getting called.
	ctx := makeRuntime(t)
	defer ctx.Finish()

	// We're in the DOWN state so EXITING shouldn't require any clean up on the supervisor/client
	ctx.runtime.Stop()

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	// In the EXITED state we expect the request to be passed on to the default router.
	ctx.defaultRouter.EXPECT().ServeHTTP(rw, req)
	ctx.runtime.ServeHTTP(rw, req)
}

func TestHandleSignal(t *testing.T) {
	ctx := makeRuntime(t)
	defer ctx.Finish()
	ctx.supervisor.EXPECT().Signal(syscall.SIGTERM)
	ctx.runtime.Signal(syscall.SIGTERM)
}

func TestHotReloadFromUPState(t *testing.T) {
	// First we'll set up a temporary directory to make changes in.
	testDir, err := ioutil.TempDir("", "hot-reload-from-up")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(testDir)

	// Now we set up the runtime and enable filesystem watching. When a change happens, we're
	// expecting to see a restart happen.
	ctx := makeRuntime(t, runtime.WatchDir(0, testDir))
	defer cleanUpWatchers(ctx)
	defer ctx.Finish()
	ctx.expectRuntimeStartTrace()
	ctx.runtime.Start()

	// When the file changes we expect to see a restart
	ctx.expectRuntimeStopTrace()
	ctx.expectRuntimeStartTrace()

	// Ok, we're set up so lets change something in the temporary directory.
	file := filepath.Join(testDir, "new-file")
	if err := ioutil.WriteFile(file, []byte{}, 0644); err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestHotReloadFromDOWNState(t *testing.T) {
	// First we'll set up a temporary directory to make changes in.
	testDir, err := ioutil.TempDir("", "hot-reload-from-down")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(testDir)

	// Now we set up the runtime and enable filesystem watching. When a change happens, we're
	// expecting to see a restart happen.
	ctx := makeRuntime(t, runtime.WatchDir(0, testDir))
	defer cleanUpWatchers(ctx)
	defer ctx.Finish()

	// As the runtime is currently down, we only expect to see a start up trace on the change.
	ctx.expectRuntimeStartTrace()

	// Ok, we're set up so lets change something in the temporary directory.
	file := filepath.Join(testDir, "new-file")
	if err := ioutil.WriteFile(file, []byte{}, 0644); err != nil {
		t.Fatalf("%+v", err)
	}
}

func cleanUpWatchers(ctx *testContext) {
	// The library we're using to watch for file changes behaves a little oddly when there are
	// multiple watchers on the go. This utility function stops the runtime, which invokes the
	// exit handler for the watcher and ensures it is cleaned up.
	// Without this the two hot reload tests work by themselves but not when the whole package
	// is run.
	ctx.expectRuntimeStopTrace()
	ctx.runtime.Stop()
}

type testContext struct {
	ctrl          *gomock.Controller
	runtime       *runtime.Runtime
	supervisor    *MockSupervisor
	client        *MockClient
	router        *MockRouter
	defaultRouter *MockRouter
	done          chan struct{}
}

func (c *testContext) Finish() {
	// Some events happen asynchronously. This sleep gives them all time to take place before
	// we close out the test and check that everything has been called correctly.
	time.Sleep(10 * time.Millisecond)

	c.ctrl.Finish()
}

var SOCKET = "/tmp/test.runtime.sock"

func (c *testContext) expectRuntimeStartTrace() {
	c.supervisor.EXPECT().Start(map[string]string{"FOLD_SOCK_ADDR": SOCKET})
	c.supervisor.EXPECT().Wait()
	c.client.EXPECT().Start(SOCKET)
	c.client.EXPECT().GetManifest(gomock.Any())
	c.router.EXPECT().Configure(gomock.Any())
}

func (c *testContext) expectRuntimeStopTrace() {
	c.client.EXPECT().Stop()
	c.supervisor.EXPECT().Stop()
}

func makeRuntime(
	t *testing.T,
	options ...runtime.Option,
) *testContext {
	ctrl := gomock.NewController(t)
	supervisor := NewMockSupervisor(ctrl)
	client := NewMockClient(ctrl)
	socketFactory := func() string { return SOCKET }
	defaultRouter := NewMockRouter(ctrl)
	activeRouter := NewMockRouter(ctrl)
	routerFactory := func(logger logging.Logger, doer router.RequestDoer) runtime.Router {
		return activeRouter
	}

	mocks := []runtime.Option{
		runtime.WithSupervisor(supervisor),
		runtime.WithClient(client),
		runtime.WithSocketFactory(socketFactory),
		runtime.WithRouterFactory(routerFactory),
		runtime.WithDefaultRouter(defaultRouter),
		runtime.OnProcessEnd(func() {}),
	}

	done := make(chan struct{})

	rt := runtime.NewRuntime(
		logging.NewTestLogger(),
		"test",
		[]string{"arg"},
		done,
		append(mocks, options...)...,
	)

	return &testContext{
		ctrl:          ctrl,
		runtime:       rt,
		supervisor:    supervisor,
		client:        client,
		router:        activeRouter,
		defaultRouter: defaultRouter,
		done:          done,
	}
}
