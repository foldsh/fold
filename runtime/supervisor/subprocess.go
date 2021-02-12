// Package subprocess offers a simple way to manage a subprocess.
// Its concern is starting, stopping, handling stdin/out... that kind of thing.
package supervisor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/foldsh/fold/logging"
)

type subprocess struct {
	logger       logging.Logger
	foldSockAddr string
	cmd          *exec.Cmd
	sout         io.Writer
	serr         io.Writer
	terminated   chan error
	health       bool
}

func newSubprocess(logger logging.Logger, foldSockAddr string) foldSubprocess {
	return &subprocess{
		logger:       logger,
		foldSockAddr: foldSockAddr,
		sout:         os.Stdout,
		serr:         os.Stderr,
		terminated:   make(chan error),
	}
}

func (sp *subprocess) run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = sp.sout
	c.Stderr = sp.serr
	c.Env = append(
		os.Environ(),
		fmt.Sprintf("FOLD_SOCK_ADDR=%s", sp.foldSockAddr),
	)
	sp.cmd = c
	sp.health = true
	go func() {
		if err := sp.cmd.Run(); err != nil {
			sp.logger.Debugf("Process has exited, determining cause")
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if exitErr.ExitCode() == -1 {
					// terminated by a signal, this is expected
					sp.logger.Debugf("The process was ended by a signal")
					sp.terminated <- nil
				} else {
					sp.logger.Debugf("The process ended unexpectedly %+v", err)
					sp.terminated <- err
				}
			} else {
				sp.logger.Debugf("The process ended unexpectedly %+v", err)
				sp.terminated <- err
			}
		} else {
			sp.terminated <- nil
			sp.health = false
		}
	}()
	return nil
}

func (sp *subprocess) wait() error {
	sp.logger.Debugf("Waiting for process to terminate")
	err := <-sp.terminated
	return err
}

func (sp *subprocess) kill() error {
	sp.logger.Debugf("Killing subprocess")
	if err := sp.cmd.Process.Kill(); err != nil {
		sp.logger.Debugf("%+v", err)
		return errors.New("failed to kill subprocess")
	}
	return nil
}

func (sp *subprocess) signal(sig os.Signal) error {
	sp.logger.Debugf("Sending signal %v to subprocess", sig)
	if sp.cmd.Process != nil {
		if err := sp.cmd.Process.Signal(sig); err != nil {
			sp.logger.Debugf("%+v", err)
			return errors.New("failed to send signal")
		}
	}
	return nil
}

func (sp *subprocess) healthz() bool {
	return sp.health
}
