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
package supervisor_test

import (
	"bytes"
	"testing"

	"github.com/foldsh/fold/runtime/supervisor"
)

func TestShouldStartAProcess(t *testing.T) {
	expectation := "TESTINPUT"
	s, sout, _ := makeProcess("echo", []string{"-n", expectation}, nil)
	if err := s.Start(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	actual := sout.String()
	if actual != expectation {
		t.Errorf("Expected %s but found %s", expectation, actual)
	}
}

func TestShouldSetEnvCorrectly(t *testing.T) {

}

func TestShouldStopAProcessGracefully(t *testing.T) {

}

func TestShouldKillAProcess(t *testing.T) {

}

func TestShouldSignalAProcess(t *testing.T) {

}

func TestShouldRestartAProcess(t *testing.T) {

}

func TestOnErrorShouldCaptureStderrAndUpdateStatus(t *testing.T) {

}

func TestInvalidCommandShouldCaptureStderrAndUpdateStatus(t *testing.T) {

}

func makeProcess(
	cmd string,
	args []string,
	env map[string]string,
) (*supervisor.Supervisor, *bytes.Buffer, *bytes.Buffer) {
	s := supervisor.NewSupervisor(cmd, args, env)
	sout := &bytes.Buffer{}
	s.Sout = sout
	serr := &bytes.Buffer{}
	s.Serr = serr
	return s, sout, serr
}
