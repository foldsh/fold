package fsm

import (
	"fmt"
	"sync"

	"github.com/foldsh/fold/logging"
)

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

type TransitionMap map[State]map[Event]Transition

type FSM struct {
	logger logging.Logger

	state      State
	stateMutex *sync.Mutex
	// This is used to keep track of all the states that have been registered with the FSM.
	states map[State]struct{}
	// This is used to keep track of all of the transitions available in each state.
	transitionMap TransitionMap
	// This keeps track of the callbacks that are to be invoked on transitioning to a given state.
	callbacks map[State][]Callback
}

func NewFSM(logger logging.Logger, initialState State, transitions Transitions) *FSM {
	fsm := &FSM{
		logger:     logger,
		state:      initialState,
		stateMutex: &sync.Mutex{},
		states:     make(map[State]struct{}),
		callbacks:  make(map[State][]Callback),
	}
	fsm.states[initialState] = struct{}{}
	transitionMap := make(TransitionMap)
	for _, t := range transitions {
		ensureEventMap(transitionMap, t)
		transitionMap[t.From][t.Event] = t
		fsm.states[t.From] = struct{}{}
		fsm.states[t.To] = struct{}{}
	}
	fsm.transitionMap = transitionMap
	return fsm
}

func (fsm *FSM) AddTransition(transition Transition) {
	ensureEventMap(fsm.transitionMap, transition)
	fsm.transitionMap[transition.From][transition.Event] = transition
}

type NoSuchStateError struct {
	State State
}

func (nsse NoSuchStateError) Error() string {
	return fmt.Sprintf("no such state '%s'", nsse.State)
}

func (fsm *FSM) OnTransitionTo(state State, callback Callback) error {
	if _, exists := fsm.states[state]; !exists {
		return NoSuchStateError{}
	}
	fsm.callbacks[state] = append(fsm.callbacks[state], callback)
	return nil
}

func ensureEventMap(tm TransitionMap, t Transition) {
	if _, exists := tm[t.From]; !exists {
		tm[t.From] = make(map[Event]Transition)
	}
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
		fsm.logger.Debugf("No such transition: event '%s' in state '%s'", event, fsm.state)
		return NoSuchTransitionError{
			From:  fsm.state,
			Event: event,
		}
	}

	fsm.logger.Debugf(
		"Event '%s': transitioning state from '%s' to '%s'",
		event,
		transition.From,
		transition.To,
	)
	// First we invoke the callbacks on the transition
	for _, cb := range transition.Callbacks {
		cb()
	}
	// Then we invoke the callbacks registered for transition to this state, if any
	if callbacks, exists := fsm.callbacks[transition.To]; exists {
		for _, cb := range callbacks {
			cb()
		}
	}
	// Finally we update the state
	fsm.state = transition.To
	return nil
}
