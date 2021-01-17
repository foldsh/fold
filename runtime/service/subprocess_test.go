package service

import (
	"bytes"
	"testing"
	"time"
)

func TestShortLivedSubprocess(t *testing.T) {
	p, err := newSubprocess(Command{"echo", []string{"-n", "TESTINPUT"}}, newAddr())
	buf := bytes.Buffer{}
	p.setStdout(&buf)
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	p.run()
	p.wait()
	out := buf.String()
	if out != "TESTINPUT" {
		t.Fatalf("Expected process to print TESTINPUT to stdout, but found %s", out)
	}
}

func TestEnvSetCorrectly(t *testing.T) {
	addr := newAddr()
	cmd := "echo -n $FOLD_SOCK_ADDR"
	p, err := newSubprocess(Command{"bash", []string{"-c", cmd}}, addr)
	buf := bytes.Buffer{}
	p.setStdout(&buf)
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	p.run()
	p.wait()
	out := buf.String()
	if out != addr {
		t.Fatalf("Expected process to print %s, but found %s", addr, out)
	}
}

func TestLongLivedProcessKill(t *testing.T) {
	// Hardly a long lived process but 10 seconds is more than enough time
	// to kill it and test that that is working.
	cmd := "echo 'for i in {1..10}; do echo $i; sleep 1; done' | bash"
	p, err := newSubprocess(Command{"bash", []string{"-c", cmd}}, newAddr())
	buf := bytes.Buffer{}
	p.setStdout(&buf)
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	p.run()
	// Ugly but we need to allow a little bit of time for the process to
	// start and print the first number.
	time.Sleep(42 * time.Millisecond)
	p.kill()
	out := buf.String()
	if out != "1\n" {
		t.Fatalf("Expected process to print 1 but found %s", out)
	}
}
