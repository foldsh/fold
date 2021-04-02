package runtime_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"testing"

	gomock "github.com/golang/mock/gomock"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime"
	"github.com/foldsh/fold/runtime/handler"
	"github.com/foldsh/fold/runtime/router"
)

var SOCKET = "/tmp/test.runtime.sock"

func TestStart(t *testing.T) {
	// Starting the runtime should start the supervisor and then the client.
	// The router should then configure itself by calling GetManifest
	rt, s, c := makeRuntime(t)
	startRuntime(rt, s, c)

	if rt.State() != runtime.UP {
		t.Errorf(
			"After succesfully starting the runtime should be in the UP state, but found %v",
			rt.State(),
		)
	}
}

func TestStop(t *testing.T) {
	rt, s, c := makeRuntime(t)
	s.EXPECT().Stop()
	c.EXPECT().Stop()
	rt.Stop()

	if rt.State() != runtime.EXITED {
		t.Errorf(
			"After stopping the runtime should be in the EXITED state, but found %v",
			rt.State(),
		)
	}
}

func TestExitOnCrash(t *testing.T) {
	// The default behaviour is to
	rt, _, _ := makeRuntime(t)
	rt.Trigger(runtime.CRASH)

	if rt.State() != runtime.EXITED {
		t.Errorf("Expect the runtime to transition to the EXITED state, but found %v", rt.State())
	}
}

func TestKeepAliveOnCrash(t *testing.T) {
	// Triggering a crash should result in the crash handler being called.
	rt, _, _ := makeRuntime(t, runtime.CrashPolicy(runtime.KEEP_ALIVE))
	rt.Trigger(runtime.CRASH)

	// TODO what methods should be called? We should ensure both the client and supervisor are
	// down and ready to start again.

	if rt.State() != runtime.DOWN {
		t.Errorf("Expecte the runtime to transition to the DOWN state, but found %v", rt.State())
	}
}

func TestDoRequestInUPState(t *testing.T) {
	// In the UP state the request should make it through to the Client
	rt, s, c := makeRuntime(t)
	startRuntime(rt, s, c)

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	rt.DoRequest(rw, req)
	// if res == nil {
	// 	t.Errorf("In the UP state there should be a response")
	// }
	// if err != nil {
	// 	t.Errorf("%+v", err)
	// }
}

func TestDoRequestInDOWNState(t *testing.T) {
	// The handler should result in the client getting called.
	rt, _, _ := makeRuntime(t)

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	rt.DoRequest(rw, req)
	// TODO check that it is the default response
	// if res == nil {
	// 	t.Errorf("In the DOWN state there should be a response")
	// }
	// if err != nil {
	// 	t.Errorf("%+v", err)
	// }
}

func TestDoRequestInEXITEDState(t *testing.T) {
	// The handler should result in the client getting called.
	rt, s, c := makeRuntime(t)

	s.EXPECT().Stop()
	c.EXPECT().Stop()
	rt.Stop()

	rw := handler.NewResponseWriter()
	req, _ := http.NewRequest("GET", "/fold", ioutil.NopCloser(strings.NewReader("fold")))
	rt.DoRequest(rw, req)
	// if res != nil {
	// 	t.Errorf("There should be no response in the EXITED state%+v", err)
	// }
	// if err == nil {
	// 	t.Errorf("DoRequest should return an error in the EXITED state.")
	// }
}

func TestHandleSignal(t *testing.T) {
	// A signal should be passed on to the supervisor.
	// TODO how to trigger the signal?
	_, s, _ := makeRuntime(t)
	s.EXPECT().Signal(syscall.SIGTERM)
}

func TestHotReloadFromUPState(t *testing.T) {
	// When enabled, a change in the file system should result in a restart
	rt, s, c := makeRuntime(t)
	startRuntime(rt, s, c)

	s.EXPECT().Stop()
	c.EXPECT().Stop()

	s.EXPECT().Start(map[string]string{"FOLD_SOCK_ADDR": SOCKET})
	c.EXPECT().Start(SOCKET)
	c.EXPECT().GetManifest(gomock.Any())

	rt.Trigger(runtime.FILE_CHANGE)
}

func TestHotReloadFromDOWNState(t *testing.T) {
	// When enabled, a change in the file system should result in a restart
	rt, s, c := makeRuntime(t)

	s.EXPECT().Start(map[string]string{"FOLD_SOCK_ADDR": SOCKET})
	c.EXPECT().Start(SOCKET)
	c.EXPECT().GetManifest(gomock.Any())

	rt.Trigger(runtime.FILE_CHANGE)
}

func makeRuntime(
	t *testing.T,
	options ...runtime.Option,
) (*runtime.Runtime, *MockSupervisor, *MockClient) {
	ctrl := gomock.NewController(t)
	supervisor := NewMockSupervisor(ctrl)
	client := NewMockClient(ctrl)
	socketFactory := func() string { return SOCKET }
	routerFactory := func(logger logging.Logger, doer router.RequestDoer) runtime.Router {
		return NewMockRouter(ctrl)
	}

	mocks := []runtime.Option{
		runtime.WithSupervisor(supervisor),
		runtime.WithClient(client),
		runtime.WithSocketFactory(socketFactory),
		runtime.WithRouterFactory(routerFactory),
	}

	rt := runtime.NewRuntime(
		logging.NewTestLogger(),
		"test",
		[]string{"arg"},
		append(mocks, options...)...,
	)

	return rt, supervisor, client
}

func startRuntime(rt *runtime.Runtime, s *MockSupervisor, c *MockClient) {
	s.EXPECT().Start(map[string]string{"FOLD_SOCK_ADDR": SOCKET})
	c.EXPECT().Start(SOCKET)
	c.EXPECT().GetManifest(gomock.Any())
	rt.Start()
}
