package fsm_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/foldsh/fold/runtime/fsm"
)

func TestBasicFSM(t *testing.T) {
	f := fsm.NewFSM(
		"off",
		fsm.Transitions{
			{"switch", "off", "on", nil},
			{"switch", "on", "off", nil},
		},
	)

	f.Emit("switch")

	if f.State() != "on" {
		t.Fatalf("State should have transitioned to 'on' but it did not")
	}

	f.Emit("switch")

	if f.State() != "off" {
		t.Fatalf("State should have transitioned to 'off' but it did not")
	}
}

func TestConcurrentTransitions(t *testing.T) {
	f := fsm.NewFSM(
		"off",
		fsm.Transitions{
			{"switch", "off", "on", nil},
			{"switch", "on", "off", nil},
		},
	)
	// We're going to create a load of goroutines that all try to transition the state.
	// If access is not synchronised then the transitions get interleaved and the state
	// deviates from what we'd expect. Note that this does not always fail without synchronisation,
	// sometimes it is correct by chance.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			f.Emit("switch")
			wg.Done()
		}()
	}
	wg.Wait()

	// We've emitted 100 switches so if everything is serialised properly we would expect to
	// end in the "off" state (any even number of switches would).
	if f.State() != "off" {
		t.Fatalf("State should have transitioned to 'off' but it did not")
	}
}

func TestNoSuchTransition(t *testing.T) {
	f := fsm.NewFSM(
		"off",
		fsm.Transitions{
			{"switch", "on", "off", nil},
		},
	)

	err := f.Emit("switch")

	var nste fsm.NoSuchTransitionError
	if !errors.As(err, &nste) {
		t.Errorf("Expected NoSuchTransitionError but found %v", err)
	}
}

type lightswitch struct {
	count    int
	onCount  int
	offCount int
	f        *fsm.FSM
}

func newLightSwitch() *lightswitch {
	l := &lightswitch{}
	f := fsm.NewFSM(
		"off",
		fsm.Transitions{
			{
				"switch",
				"off",
				"on",
				[]fsm.Callback{func() { l.count++ }, func() { l.onCount++ }},
			},
			{
				"switch",
				"on",
				"off",
				[]fsm.Callback{func() { l.count++ }, func() { l.offCount++ }},
			},
		},
	)
	l.f = f
	return l
}

func (l *lightswitch) pressSwitch() {
	l.f.Emit("switch")
}

func TestCallbacks(t *testing.T) {
	l := newLightSwitch()

	for i := 0; i < 100; i++ {
		l.pressSwitch()
	}

	if l.count != 100 {
		t.Errorf("Expected 100 presses but found %d", l.count)
	}
	if l.onCount != 50 {
		t.Errorf("Expected 100 on presses but found %d", l.onCount)
	}
	if l.offCount != 50 {
		t.Errorf("Expected 100 off presses but found %d", l.offCount)
	}
}
