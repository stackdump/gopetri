package statemachine_test

import (
	"fmt"
	"github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterMachine(t *testing.T) {
	s := statemachine.LoadPnmlFromFile("../examples/counter.xml")

	t.Run("load machine ", func(t *testing.T) {
		s.Init()
		if len(s.Transitions) == 0 {
			t.Fatal("failed to load statemachine")
		}
	})

	t.Run("commit actions", func(t *testing.T) {
		commit := func(action string, expectFail bool) {
			res, err := s.Commit(action, 1)
			fmt.Printf("output: %v expectFail: %v\n", res, expectFail)

			if expectFail && err == nil {
				t.Fatalf("expected %v to fail", action)
			}

			if !expectFail && err != nil {
				t.Fatalf("expected %v not to fail", action)
			}
		}

		commit("DEC_0", true)
		commit("INC_0", false)
		commit("INC_0", false)
	})
}

// NOTE: if you do such a test on an unbounded network (like the one above)
// be warned that golang will happily try to recurse forever
func walkNet(sm statemachine.StateMachine, games *uint64) {
	actions, ok := sm.ValidActions(1)
	if !ok {
		*games++
	}

	for _, state := range actions {
		walkNet(sm.Clone(state), games)
	}
}

const nineFactorial uint64 = 362880

// the actual test for boundedness is that
// the recursive function walkNet should not infinitely recurse
// this test completes in ~7s
func TestTicTacToeStateSpace(t *testing.T) {
	s := statemachine.LoadPnmlFromFile("../examples/octoe.xml")
	// remove extraneous early game-ending actions
	delete(s.Transitions, statemachine.Action("END_O"))
	delete(s.Transitions, statemachine.Action("END_X"))

	t.Run("initialize game", func(t *testing.T){
		s.Init()
		_, err := s.Commit("EXEC", 1)
		assert.Nil(t, err)
	})

	t.Run("walk state space", func(t *testing.T) {
		var games uint64 = 0
		walkNet(s, &games)
		fmt.Printf("games: %v", games)
		assert.Equal(t, games, nineFactorial, "expected count to all possible games")
	})
}
