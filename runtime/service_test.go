/**
TODO
- test starting
- test restart on crash
- test log on crash
- test handling request
- test handling signal
- test restart on file change
*/
package runtime_test

func TestStart(t *testing.T) {
}

func TestRestartOnCrash(t *testing.T) {
}

func TestLogOnCrash(t *testing.T) {
}

func TestHandleRequest(t *testing.T) {
	// The handler should result in the client getting called.
}

func TestHandleRequestsDuringRestart(t *testing.T) {
	// If a request comes in during a restart we should wait until the restart is complete
	// in order to serve the request.
}

func TestHandleRequestsWhenCrashed(t *testing.T) {
	// If a request comes in when we're crashed, we should give a 500 response.
}

func TestHandleSignal(t *testing.T) {
	// A signal should be passed on to the supervisor.
}

func TestHotReload(t *testing.T) {
	// When enabled, a change in the file system should result in a restart
}
