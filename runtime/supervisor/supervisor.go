package supervisor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/foldsh/fold/logging"
)

var (
	TerminatedBySignal = errors.New("process terminated by a signal")
)

type Status int

const (
	NOTSTARTED Status = iota + 1
	STARTFAILED
	RUNNING
	CRASHED
	COMPLETE
)

type Supervisor struct {
	Cmd        string
	Args       []string
	Env        map[string]string
	Sout       io.Writer
	Serr       io.Writer
	Terminated chan error

	state   Status
	command *exec.Cmd
	logger  logging.Logger
}

func NewSupervisor(
	logger logging.Logger,
	cmd string,
	args []string,
	env map[string]string,
) *Supervisor {
	return &Supervisor{
		Cmd:        cmd,
		Args:       args,
		Env:        env,
		Sout:       os.Stdout,
		Serr:       os.Stderr,
		Terminated: make(chan error, 1),
		logger:     logger,
		state:      NOTSTARTED,
	}
}

func (s *Supervisor) Status() Status {
	return s.state
}

func (s *Supervisor) Start() error {
	command := exec.Command(s.Cmd, s.Args...)
	s.command = command
	command.Stdout = s.Sout
	command.Stderr = s.Serr
	command.Env = os.Environ()
	for key, value := range s.Env {
		command.Env = append(
			command.Env,
			fmt.Sprintf("%s=%s", key, value),
		)
	}
	err := command.Start()
	if err != nil {
		return fmt.Errorf("process failed to start: %v", err)
	}
	s.state = RUNNING
	go func() {
		err := command.Wait()
		if err == nil {
			// The command executed successfully so there is nothing left to do.
			s.state = COMPLETE
			s.Terminated <- nil
			return
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			// Can this even happen given that we have used Start?
			s.state = STARTFAILED
			s.Terminated <- fmt.Errorf("process did not run successfully: %v", err)
			return
		}
		// It's an exit error, so the process ran but stopped for some reason.
		if exitErr.ExitCode() == -1 {
			// Terminated by a signal, so this is expected.
			s.logger.Debugf("The process was terminated by a signal")
			s.state = COMPLETE
			s.Terminated <- TerminatedBySignal
			return
		}
		// The users program crashed, they have a bug.
		s.logger.Debugf("The process ended unexpectedly %+v", err)
		s.state = CRASHED
		s.Terminated <- fmt.Errorf("process crashed: %v", err)
		return
	}()
	return nil
}

func (s *Supervisor) Restart() error {
	return nil
}

func (s *Supervisor) Stop() error {
	return nil
}

func (s *Supervisor) Kill() error {
	return s.command.Process.Kill()
}

func (s *Supervisor) Wait() error {
	err := <-s.Terminated
	return err
}

func (s *Supervisor) Signal(sig os.Signal) error {
	return s.command.Process.Signal(sig)
}
