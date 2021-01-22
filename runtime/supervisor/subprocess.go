// Package subprocess offers a simple way to manage a subprocess.
// Its concern is starting, stopping, handling stdin/out... that kind of thing.
package supervisor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type subprocess struct {
	foldSockAddr string
	cmd          *exec.Cmd
	sout         io.Writer
	serr         io.Writer
}

/*
We want to be able to hide this behind something so that we can completely
swap the subprocess out for the threaded version in a test.

why can't we d
*/

func newSubprocess(foldSockAddr string) *subprocess {
	return &subprocess{
		foldSockAddr: foldSockAddr,
		sout:         os.Stdout,
		serr:         os.Stderr,
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
