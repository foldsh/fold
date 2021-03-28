package supervisor_test

import (
	"bytes"
	"errors"
	"syscall"
	"testing"

	"github.com/foldsh/fold/logging"
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
	status := s.Status()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldSetEnvCorrectly(t *testing.T) {
	cmd := "echo -n $ONE && echo -n $TWO"
	s, sout, _ := makeProcess(
		"bash",
		[]string{"-c", cmd},
		map[string]string{"ONE": "ONE", "TWO": "TWO"},
	)
	if err := s.Start(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	actual := sout.String()
	expectation := "ONETWO"
	if actual != expectation {
		t.Errorf("Expected %s but found %s", expectation, actual)
	}
	status := s.Status()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldStopAProcessGracefully(t *testing.T) {

}

func TestShouldKillAProcess(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"}, nil)
	if err := s.Start(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Kill(); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	status := s.Status()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldSignalAProcess(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"}, nil)
	if err := s.Start(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	status := s.Status()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
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
	s := supervisor.NewSupervisor(logging.NewTestLogger(), cmd, args, env)
	sout := &bytes.Buffer{}
	s.Sout = sout
	serr := &bytes.Buffer{}
	s.Serr = serr
	return s, sout, serr
}
