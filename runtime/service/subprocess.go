// Package subprocess offers a simple way to manage a subprocess.
// Its concern is starting, stopping, handling stdin/out... that kind of thing.
package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type subprocess struct {
	cmd *exec.Cmd
}

func newSubprocess(cmd Command, foldSockAddr string) (*subprocess, error) {
	c := exec.Command(cmd.Command, cmd.Args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(), fmt.Sprintf("FOLD_SOCK_ADDR=%s", foldSockAddr))
	return &subprocess{c}, nil
}

func (sp *subprocess) run() error {
	if err := sp.cmd.Start(); err != nil {
		return errors.New("failed to start subprocess")
	}
	return nil
}

func (sp *subprocess) wait() error {
	if err := sp.cmd.Wait(); err != nil {
		return errors.New("failed to wait for subprocess")
	}
	return nil
}

func (sp *subprocess) kill() error {
	if err := sp.cmd.Process.Kill(); err != nil {
		return errors.New("failed to kill subprocess")
	}
	return nil
}

// TODO handle shut down on signal: Just SIGTERM is fine for now.
func (sp *subprocess) signal(sig os.Signal) error {
	if err := sp.cmd.Process.Signal(sig); err != nil {
		return errors.New("failed to send signal")
	}
	return nil
}

// Useful for capturing stdout for testing.
func (sp *subprocess) setStdout(w io.Writer) {
	sp.cmd.Stdout = w
}

// Useful for capturing stderr for testing.
func (sp *subprocess) setStderr(w io.Writer) {
	sp.cmd.Stderr = w
}
