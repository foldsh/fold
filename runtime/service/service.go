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
	"errors"
	"os"

	"github.com/foldsh/fold/manifest"
)

type Service interface {
	Start()
	DoRequest(req *Request) *Response
	GetManifest() *manifest.Manifest
	Signal(sig os.Signal)
}

func NewService(cmd Command) (Service, error) {
	addr := newAddr()
	client := newIngressClient(addr)
	process, err := newSubprocess(cmd, addr)
	if err != nil {
		return nil, errors.New("Failed to start service.")
	}
	return &service{cmd, client, process}, nil
}

type Command struct {
	command string
	args    []string
}

type service struct {
	cmd     Command
	client  *ingressClient
	process *subprocess
}

func (s *service) Start() error {
	// TODO startup the process, and then the client.
	//  There is one challenge here, which is that until the server
	//  boots up on the other side, the socket doesn't exist.
	//  this obviously leads to the client failing miserably.
	err := s.process.run()
	if err != nil {
		return err
	}
	err = s.client.start()
	if err != nil {
		return errors.New("Failed to start client.")
	}
	return nil
}

func (s *service) DoRequest(req *Request) *Response {
	return &Response{}
}

func (s *service) GetManifest() *manifest.Manifest {
	return &manifest.Manifest{}
}

func (s *service) Signal(sig os.Signal) {
}
