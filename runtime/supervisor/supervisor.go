// Package service is responsible for the management of a single service, defined
// by the user. It does this through the management of a subprocess, which is running
// an application defined with the fold sdk.
//
// It exposes a simple RPC based interface for interacting with that subprocess
// which can be used by the handlers as they wish. Communication between the
// runtime and the subprocess is done via gRPC over a unix domain socket.
//
// This approach was chosen because it results in very little overhead (it adds
// a few milliseconds to startup time and a few hundred microseconds to each
// request) and because it allows us to generate much of the code for both sides.
// This will make it very easy to implement the sdk in multiple languages.
package supervisor

import (
	"context"
	"errors"
	"os"
	"syscall"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

var (
	FailedToStartProcess  = errors.New("failed to start subprocess")
	FailedToStopProcess   = errors.New("failed to stop subprocess")
	FailedToSignalProcess = errors.New("failed to signal subprocess")
	CannotServiceRequest  = []byte(
		`{"title": "Cannot service request","detail":"The application is not running. Please check the logs."}`,
	)
)

// This is just a wrapper around the protobuf definition in proto/ingress.proto
// It makes them easier to use and avoids exposing the generated code to the
// rest of the runtime package.
type Request struct {
	HttpMethod  string
	Handler     string
	Path        string
	Body        []byte
	Headers     map[string][]string
	PathParams  map[string]string
	QueryParams map[string][]string
}

// This is just a wrapper around the protobuf definition in proto/ingress.proto
// It makes them easier to use and avoids exposing the generated code to the
// rest of the runtime package.
type Response struct {
	Status  int
	Body    []byte
	Headers map[string][]string
}

func NewSupervisor(logger logging.Logger, cmd string, args ...string) *Supervisor {
	service := &Supervisor{
		cmd:               cmd,
		args:              args,
		logger:            logger,
		clientFactory:     newIngressClient,
		subprocessFactory: newSubprocess,
	}
	return service
}

// This struct is basically about managing all the different states for a client and subprocess.
// States -> Action:
//   - process: up, client: up -> Do Request
//   - process: up, client: down -> Error Response
//   - process: down, client: up -> Error Response
//   - process: down, client: down -> Error Response
// This essentially boils down to two states - Available and NotAvailable.
// The Available state can only be reached by succesfully calling the Start method - which
// means that both the process and client have been started succesfully.
// In all states the process should remain available for restart.
// Setting the state will need to be synchronised as this is shared by the router which can
// potentially issue multiple request at once.
type Supervisor struct {
	addr              string
	cmd               string
	args              []string
	logger            logging.Logger
	clientFactory     clientFactory
	client            *ingressClient
	subprocessFactory subprocessFactory
	process           foldSubprocess
}

type clientFactory func(logging.Logger, string) *ingressClient
type subprocessFactory func(logging.Logger, string) foldSubprocess

type foldSubprocess interface {
	run(string, ...string) error
	kill() error
	// wait() error
	signal(os.Signal) error
	healthz() bool
}

func (s *Supervisor) Start() error {
	s.addr = newAddr()
	s.logger.Debugf("starting subprocess with socket %s", s.addr)
	s.process = s.subprocessFactory(s.logger, s.addr)
	if err := s.process.run(s.cmd, s.args...); err != nil {
		return FailedToStartProcess
	}

	s.client = s.clientFactory(s.logger, s.addr)
	s.logger.Debugf("starting gRPC client")
	if err := s.client.start(); err != nil {
		s.logger.Debugf("failed to start gRPC client")
		return FailedToStartProcess
	}
	return nil
}

func (s *Supervisor) Stop() error {
	if err := s.Signal(syscall.SIGTERM); err != nil {
		s.logger.Debugf("failed to stop subprocess")
		return FailedToStopProcess
	}
	return nil
}

func (s *Supervisor) Restart() error {
	s.logger.Debugf("restarting subprocess")
	if s.process.healthz() {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	if err := s.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Supervisor) DoRequest(req *Request) (*Response, error) {
	s.logger.Debug("performing application request: ", req)
	if s.process.healthz() {
		return s.client.doRequest(context.Background(), req)
	} else {
		// Bit of a weird error state, but this is quite an easy way to accomplish
		// the desired effect. Basically, the supervisor should remain responsive when
		// the child process is unavailable. I.e., just because the subprocess is unhealthy,
		// it doesn't mean the runtime is too.
		res := &Response{Status: 502, Body: CannotServiceRequest}
		return res, nil
	}
}

func (s *Supervisor) GetManifest() (*manifest.Manifest, error) {
	s.logger.Debugf("loading application manifest")
	return s.client.getManifest(context.Background())
}

func (s *Supervisor) Signal(sig os.Signal) error {
	s.logger.Debug("supervisor received signal ", sig)
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		s.logger.Debugf("shutting down application subprocess")
		defer os.Remove(s.addr)
		if err := s.process.signal(sig); err != nil {
			s.logger.Debugf("failed to signal subprocess %+v", err)
			return FailedToStopProcess
		}
		// if err := s.process.wait(); err != nil {
		// 	s.logger.Debugf("failed to wait for subprocess %+v", err)
		// 	return FailedToStopProcess
		// }
	case syscall.SIGKILL:
		s.logger.Debugf("killing the application subprocess")
		defer os.Remove(s.addr)
		if err := s.process.kill(); err != nil {
			s.logger.Debugf("failed to kill subprocess")
			return FailedToStopProcess
		}
	}
	return nil
}
