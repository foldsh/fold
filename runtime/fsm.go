package runtime

import "sync"

type RuntimeState uint8

const (
	DOWN RuntimeState = iota + 1
	UP
	EXITED
)

type EventT uint8

const (
	START EventT = iota + 1
	STOP
	CRASH
	FILE_CHANGE
)

type EventHandler func() error

type fsm struct {
	stateMu *sync.Mutex
}

func newfsm() *fsm {
}

func (f *fsm) Transition() {
}

func (f *fsm) Trigger() {
}
