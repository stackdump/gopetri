// Packaget ptnet provides a place-transition equivilent of an elementary petri-net
package statemachine

import (
	"errors"
)

type StateVector []uint64
type Transition []int64
type Action string

type StateMachine struct {
	Initial	StateVector
	Transitions map[Action]Transition
	State	StateVector
}

// initial state - this value also
// serves as a bounds check - by convention no entry in state vector
// should ever exceed the initial value
func (s StateMachine) Init(){
	if len(s.State) == 0 {
		for offset, val := range s.Initial {
			s.State[offset] = val
		}
	}
}

// test that state has not exceeded initial values
func (s StateMachine) InBounds() bool {
	for offset, val := range s.Initial {
		if s.State[offset] > val {
			return false
		}
	}
	return true
}

// apply the transformation without overwriting state
func (s StateMachine) Transform(transform Transition, multiplier uint64) ([]int64, error) {
	var vectorOut []int64
	var err error = nil

	for offset, delta := range transform {
		val := int64(s.State[offset]) + delta*int64(multiplier)
		vectorOut = append(vectorOut, val)
		if err != nil && val < 0 {
			err = errors.New("invalid output")
		}
	}
	return vectorOut, err
}

// apply the transformation and overwrite state
func (s StateMachine) Commit(transform Transition, multiplier uint64) ([]int64, error) {
	vectorOut, err := s.Transform(transform, multiplier)

	if err == nil{
		for offset, val := range vectorOut {
			s.State[offset] = uint64(val)
		}
	}

	return vectorOut, err
}
