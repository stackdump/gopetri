// Packaget ptnet provides a place-transition equivilent of an elementary petri-net
package statemachine

import (
	"errors"
)

type StateVector []uint64
type Transition []int64
type Action string

type StateMachine struct {
	Initial     StateVector
	Capacity    StateVector
	Transitions map[Action]Transition
	State       StateVector
}

func (s *StateMachine) Init() {
	for _, val := range s.Initial {
		s.State = append(s.State, val)
	}
}

// test that state has not exceeded initial values
func (s StateMachine) InBounds() bool {
	for offset, val := range s.Capacity {
		if s.State[offset] > val {
			return false
		}
	}
	return true
}

// apply the transformation without overwriting state
func (s StateMachine) Transform(action string, multiplier uint64) ([]int64, error) {
	var vectorOut []int64
	var err error = nil

	for offset, delta := range s.Transitions[Action(action)] {
		val := int64(s.State[offset]) + delta*int64(multiplier)
		vectorOut = append(vectorOut, val)
		if err == nil && val < 0 {
			err = errors.New("invalid output")
		}
		if err == nil && s.Capacity[offset] != 0 && val > int64(s.Capacity[offset]) {
			err = errors.New("exceeded capacity")
		}
	}
	return vectorOut, err
}

// apply the transformation and overwrite state
func (s *StateMachine) Commit(action string, multiplier uint64) ([]int64, error) {
	vectorOut, err := s.Transform(action, multiplier)

	if err == nil {
		for offset, val := range vectorOut {
			s.State[offset] = uint64(val)
		}
	}

	return vectorOut, err
}
