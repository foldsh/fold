/**
TODO
	- start a process
	- restart a process
	- stop a process
	- wait for it to finish
	- signal the process
	- expose stdout
	- expose stderr
	- model process state UP / DOWN
*/
package supervisor

import (
	"io"
	"os"

	"github.com/foldsh/fold/logging"
)

type Status int

const (
	UP Status = iota + 1
	DOWN
)

type Supervisor struct {
	Cmd  string
	Args []string
	Sout io.Writer
	Serr io.Writer

	state  Status
	logger logging.Logger
}

func NewSupervisor(cmd string, args []string, env map[string]string) *Supervisor {
	return &Supervisor{Cmd: cmd, Args: args, Sout: os.Stdout, Serr: os.Stderr}
}

func (s *Supervisor) Status() Status {
	return s.state
}

func (s *Supervisor) Start() error {
	return nil
}

func (s *Supervisor) Restart() error {
	return nil
}

func (s *Supervisor) Stop() error {
	return nil
}

func (s *Supervisor) Kill() error {
	return nil
}

func (s *Supervisor) Wait() error {
	return nil
}

func (s *Supervisor) Signal(sig os.Signal) error {
	return nil
}
