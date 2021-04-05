package supervisor_test

import (
	"bytes"
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/foldsh/fold/logging"
	"github.com/foldsh/fold/runtime/supervisor"
)

func TestShouldStartAProcess(t *testing.T) {
	expectation := "TESTINPUT"
	s, sout, _ := makeProcess("echo", []string{"-n", expectation})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	status := s.State()
	if status != supervisor.RUNNING {
		t.Errorf("Expected RUNNING but found %v", status)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	actual := sout.String()
	if actual != expectation {
		t.Errorf("Expected %s but found %s", expectation, actual)
	}
	status = s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldSetEnvCorrectly(t *testing.T) {
	s, sout, _ := makeProcess(
		"bash",
		[]string{"./testdata/env.sh"},
	)
	if err := s.Start(map[string]string{"ONE": "ONE", "TWO": "TWO"}); err != nil {
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
	status := s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldStopAProcessGracefully(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	status := s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldKillAProcess(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Kill(); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	status := s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldSignalAProcess(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	status := s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestShouldRestartAProcess(t *testing.T) {
	s, sout, _ := makeProcess(
		"bash",
		[]string{"./testdata/restart.sh"},
	)
	if err := s.Start(nil); err != nil {
		t.Fatalf("%+v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := s.Restart(nil); err != nil {
		t.Fatalf("%+v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := s.Stop(); err != nil {
		t.Fatalf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	actual := sout.String()
	expectation := "FOLDFOLD"
	if actual != expectation {
		t.Errorf("Expected %s but found %s", expectation, actual)
	}
	status := s.State()
	if status != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", status)
	}
}

func TestOnErrorShouldCaptureStderrAndUpdateStatus(t *testing.T) {
	s, sout, serr := makeProcess(
		"bash",
		[]string{"./testdata/error.sh"},
	)
	if err := s.Start(nil); err != nil {
		t.Fatalf("%+v", err)
	}
	err := s.Wait()
	var pe supervisor.ProcessError
	if !errors.As(err, &pe) {
		t.Errorf("Expected ProcessError but found %v", err)
	}
	outResult := sout.String()
	if outResult != "" {
		t.Errorf("Expected an empty string but found %s", outResult)
	}
	errResult := serr.String()
	if errResult != "expr: division by zero\n" {
		t.Errorf("Expected 'expr: division by zero' but found %s", errResult)
	}
	if s.State() != supervisor.CRASHED {
		t.Errorf("Expected CRASHED but found %v", s.State())
	}
}

func TestInvalidCommandShouldErrorAndUpdateStatus(t *testing.T) {
	s, _, _ := makeProcess("not-a-command", []string{})
	err := s.Start(nil)
	var pe supervisor.ProcessError
	if !errors.As(err, &pe) {
		t.Errorf("Expected ProcessError but found %v", err)
	}
	if s.State() != supervisor.STARTFAILED {
		t.Errorf("Expected STARTFAILED but found %v", s.State())
	}
}

func TestStopShouldBeIdempotent(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	// Ok we've stopped the process. Calling Stop and Wait again should both just complete
	// immediately without error.
	if err := s.Stop(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	if s.State() != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", s.State())
	}
}

func TestKillShouldBeIdempotent(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Kill(); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	// Ok we've stopped the process. Calling Stop and Wait again should both just complete
	// immediately without error.
	if err := s.Kill(); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	if s.State() != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", s.State())
	}
}

func TestSignalShouldBeIdempotent(t *testing.T) {
	s, _, _ := makeProcess("sleep", []string{"999"})
	if err := s.Start(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("%+v", err)
	}
	err := s.Wait()
	if !errors.Is(err, supervisor.TerminatedBySignal) {
		t.Errorf("Expected TerminatedBySignal but found %+v", err)
	}
	// Ok we've stopped the process. Calling Stop and Wait again should both just complete
	// immediately without error.
	if err := s.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("%+v", err)
	}
	if err := s.Wait(); err != nil {
		t.Errorf("%+v", err)
	}
	if s.State() != supervisor.COMPLETE {
		t.Errorf("Expected COMPLETE but found %v", s.State())
	}
}

func makeProcess(
	cmd string,
	args []string,
) (*supervisor.Supervisor, *bytes.Buffer, *bytes.Buffer) {
	sout := &bytes.Buffer{}
	serr := &bytes.Buffer{}
	s := supervisor.NewSupervisor(logging.NewTestLogger(), cmd, args, sout, serr)
	return s, sout, serr
}
