package fsm

import (
	"fmt"
	"sync"
)

/**
We need to:
 - say which events are available in a given state
 - say which transitions are available
 - expose hooks to update state on the host struct on transitions
 - synchronise state changes
*/

type State string

type Event string

type Callback func()

type Transition struct {
	Event     Event
	From      State
	To        State
	Callbacks []Callback
}

type Transitions []Transition

type FSM struct {
	state         State
	stateMutex    *sync.Mutex
	transitionMap map[State]map[Event]Transition
}

func NewFSM(initialState State, transitions Transitions) *FSM {
	fsm := &FSM{
		state:      initialState,
		stateMutex: &sync.Mutex{},
	}
	transitionMap := make(map[State]map[Event]Transition)
	for _, t := range transitions {
		if _, exists := transitionMap[t.From]; !exists {
			transitionMap[t.From] = make(map[Event]Transition)
		}
		transitionMap[t.From][t.Event] = t
	}
	fsm.transitionMap = transitionMap
	return fsm
}

func (fsm *FSM) State() State {
	fsm.stateMutex.Lock()
	defer fsm.stateMutex.Unlock()
	return fsm.state
}

type NoSuchTransitionError struct {
	From  State
	Event Event
}

func (nste NoSuchTransitionError) Error() string {
	return fmt.Sprintf(
		"no transition from state %s for event %s",
		nste.From,
		nste.Event,
	)
}

func (fsm *FSM) Emit(event Event) error {
	fsm.stateMutex.Lock()
	defer fsm.stateMutex.Unlock()

	transitions := fsm.transitionMap[fsm.state]

	transition, exists := transitions[event]
	if !exists {
		return NoSuchTransitionError{
			From:  fsm.state,
			Event: event,
		}
	}
	fsm.state = transition.To
	for _, cb := range transition.Callbacks {
		cb()
	}
	return nil
}
