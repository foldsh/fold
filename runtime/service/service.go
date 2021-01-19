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
package service

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/manifest"
)

type Service interface {
	Start() error
	Stop() error
	DoRequest(req *Request) (*Response, error)
	GetManifest() (*manifest.Manifest, error)
	Signal(sig os.Signal)
}

func NewService(logger logging.Logger, cmd Command) (Service, error) {
	addr := newAddr()
	client := newIngressClient(addr)
	process, err := newSubprocess(cmd, addr)
	if err != nil {
		return nil, errors.New("failed to start service")
	}
	return &service{cmd, addr, client, process, logger}, nil
}

type Command struct {
	Command string
	Args    []string
}

type service struct {
	cmd     Command
	addr    string
	client  *ingressClient
	process foldSubprocess
	logger  logging.Logger
}

type foldSubprocess interface {
	run() error
	wait() error
	kill() error
	signal(sig os.Signal) error
	setStdout(w io.Writer)
	setStderr(w io.Writer)
}

func (s *service) Start() error {
	s.logger.Debugf("starting application subprocess")
	// TODO startup the process, and then the client.
	//  There is one challenge here, which is that until the server
	//  boots up on the other side, the socket doesn't exist.
	//  this obviously leads to the client failing miserably.
	err := s.process.run()
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

func (s *service) Stop() error {
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
