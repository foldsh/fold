package supervisor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/foldsh/fold/logging"
)

type State int

const (
	NOTSTARTED State = iota + 1
	STARTFAILED
	RUNNING
	CRASHED
	COMPLETE
)

type Supervisor struct {
	Cmd        string
	Args       []string
	Sout       io.Writer
	Serr       io.Writer
	Terminated chan error

	state      State
	stateMutex *sync.Mutex
	command    *exec.Cmd
	logger     logging.Logger
}

func NewSupervisor(
	logger logging.Logger,
	cmd string,
	args []string,
	sout io.Writer,
	serr io.Writer,
) *Supervisor {
	return &Supervisor{
		Cmd:        cmd,
		Args:       args,
		Sout:       sout,
		Serr:       serr,
		Terminated: make(chan error, 1),
		logger:     logger,
		state:      NOTSTARTED,
		stateMutex: &sync.Mutex{},
	}
}

func (s *Supervisor) State() State {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()
	return s.state
}

func (s *Supervisor) setState(state State) {
	s.stateMutex.Lock()
	defer s.stateMutex.Unlock()
	s.state = state
}

var (
	TerminatedBySignal = errors.New("process terminated by a signal")
)

type ProcessError struct {
	Reason string
	Inner  error
}

func (e ProcessError) Error() string {
	return fmt.Sprintf("%s: %s", e.Reason, e.Inner.Error())
}

func (e ProcessError) Unwrap() error {
	return e.Inner
}

func (s *Supervisor) Start(env map[string]string) error {
	s.logger.Debugf("Starting the process")
	command := exec.Command(s.Cmd, s.Args...)
	s.command = command
	command.Stdout = s.Sout
	command.Stderr = s.Serr
	command.Env = os.Environ()
	for key, value := range env {
		command.Env = append(
			command.Env,
			fmt.Sprintf("%s=%s", key, value),
		)
	}
	err := command.Start()
	if err != nil {
		s.setState(STARTFAILED)
		return ProcessError{Reason: "process failed to start", Inner: err}
	}
	s.setState(RUNNING)
	go func() {
		err := command.Wait()
		if err == nil {
			// The command executed successfully so there is nothing left to do.
			s.setState(COMPLETE)
			s.Terminated <- nil
			return
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			// Can this even happen given that we have used Start?
			s.setState(STARTFAILED)
			s.Terminated <- ProcessError{Reason: "process did not run successfully", Inner: err}
			return
		}
		// It's an exit error, so the process ran but stopped for some reason.
		if exitErr.ExitCode() == -1 {
			// Terminated by a signal, so this is expected.
			s.logger.Debugf("The process was terminated by a signal")
			s.setState(COMPLETE)
			s.Terminated <- TerminatedBySignal
			return
		}
		// The users program crashed, they have a bug.
		s.logger.Debugf("The process ended unexpectedly %+v", err)
		s.setState(CRASHED)
		s.Terminated <- ProcessError{Reason: "process crashed", Inner: err}
		return
	}()
	return nil
}

func (s *Supervisor) Restart(env map[string]string) error {
	s.logger.Debugf("Restarting the process")
	if err := s.Stop(); err != nil {
		return fmt.Errorf("failed to restart process: %v", err)
	}
	if err := s.Wait(); err != nil && !errors.Is(err, TerminatedBySignal) {
		return fmt.Errorf("failed to restart process: %v", err)
	}
	return s.Start(env)
}

func (s *Supervisor) Stop() error {
	s.logger.Debugf("Stopping the process")
	if s.State() != RUNNING {
		return nil
	}
	return s.Signal(syscall.SIGTERM)
}

func (s *Supervisor) Kill() error {
	s.logger.Debugf("Killing the process")
	if s.State() != RUNNING {
		return nil
	}
	return s.command.Process.Kill()
}

func (s *Supervisor) Wait() error {
	s.logger.Debugf("Waiting for the process to terminate")
	if s.State() != RUNNING {
		return nil
	}
	return <-s.Terminated
}

func (s *Supervisor) Signal(sig os.Signal) error {
	if s.State() != RUNNING {
		return nil
	}
	return s.command.Process.Signal(sig)
}
