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

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
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

type Supervisor interface {
	Exec(string, ...string) error
	Shutdown() error
	DoRequest(*Request) (*Response, error)
	GetManifest() (*manifest.Manifest, error)
	Signal(os.Signal)
}

func NewSupervisor(logger logging.Logger) Supervisor {
	addr := newAddr()
	service := &service{
		addr:    addr,
		client:  newIngressClient(addr),
		process: newSubprocess(addr),
		logger:  logger,
	}
	return service
}

type service struct {
	addr    string
	client  *ingressClient
	process foldSubprocess
	logger  logging.Logger
}

type foldSubprocess interface {
	run(string, ...string) error
	wait() error
	kill() error
	signal(os.Signal) error
}

func (s *service) Exec(cmd string, args ...string) error {
	s.logger.Debugf("starting application subprocess")
	err := s.process.run(cmd, args...)
	if err != nil {
		return err
	}
	s.logger.Debugf("starting gRPC client")
	err = s.client.start()
	if err != nil {
		return errors.New("failed to start client")
	}
	return nil
}

func (s *service) Shutdown() error {
	s.logger.Debugf("shutting down application subprocess")
	os.Remove(s.addr)
	return s.process.kill()
}

func (s *service) DoRequest(req *Request) (*Response, error) {
	s.logger.Debug("performing application request: ", req)
	return s.client.doRequest(context.Background(), req)
}

func (s *service) GetManifest() (*manifest.Manifest, error) {
	s.logger.Debugf("retrieving manifest")
	return s.client.getManifest(context.Background())
}

func (s *service) Signal(sig os.Signal) {
}
