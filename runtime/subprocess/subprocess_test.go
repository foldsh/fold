package supervisor

import (
	"bytes"
	"testing"
	"time"

	"github.com/foldsh/fold/logging"
)

func TestShortLivedSubprocess(t *testing.T) {
	buf := &bytes.Buffer{}
	p := testSubProcess(buf)
	err := p.run("echo", "-n", "TESTINPUT")
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	if err := p.wait(); err != nil {
		t.Fatalf("Expected process to end without error")
	}
	out := buf.String()
	if out != "TESTINPUT" {
		t.Errorf("Expected process to print TESTINPUT to stdout, but found %s", out)
	}
}

func TestEnvSetCorrectly(t *testing.T) {
	buf := &bytes.Buffer{}
	p := testSubProcess(buf)
	cmd := "echo -n $FOLD_SOCK_ADDR"
	err := p.run("bash", "-c", cmd)
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	err = p.wait()
	if err != nil {
		t.Fatalf("Expected process to end without error")
	}
	out := buf.String()
	if out != p.foldSockAddr {
		t.Errorf("Expected process to print %s, but found %s", p.foldSockAddr, out)
	}
}

func TestLongLivedProcessKill(t *testing.T) {
	// Hardly a long lived process but 10 seconds is more than enough time
	// to kill it and test that that is working.
	buf := &bytes.Buffer{}
	p := testSubProcess(buf)
	cmd := "echo 'for i in {1..10}; do echo $i; sleep 1; done' | bash"
	err := p.run("bash", "-c", cmd)
	if err != nil {
		t.Fatalf("Failed to start subprocess")
	}
	// Ugly but we need to allow a little bit of time for the process to
	// start and print the first number.
	time.Sleep(42 * time.Millisecond)
	p.kill()
	out := buf.String()
	if out != "1\n" {
		t.Fatalf("Expected process to print 1 but found %s", out)
	}
}

func testSubProcess(b *bytes.Buffer) *subprocess {
	p := &subprocess{
		logger:       logging.NewTestLogger(),
		foldSockAddr: newAddr(),
		terminated:   make(chan error),
	}
	p.sout = b
	p.logger.Debugf("Created test process %+v", p)
	return p
}
