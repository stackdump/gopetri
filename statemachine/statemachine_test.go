package statemachine_test

import (
	"fmt"
	"github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterMachine(t *testing.T) {
	_, s := statemachine.LoadPnmlFromFile("../examples/counter.xml")

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
		commit("DEC_0", false)
		commit("INC_0", false)
	})
}

var winningSets []map[string]bool = []map[string]bool{
	{"00": true, "01": true, "02": true},
	{"10": true, "11": true, "12": true},
	{"20": true, "21": true, "22": true},
	{"00": true, "11": true, "22": true},
	{"02": true, "11": true, "20": true},
	{"00": true, "10": true, "20": true},
	{"01": true, "11": true, "21": true},
	{"02": true, "12": true, "22": true},
}

var HistoryIndex [][]statemachine.Action = [][]statemachine.Action{}

type StateMachineWithHistory struct {
	statemachine.StateMachine
	History []statemachine.Action
}

func containsWin(history map[string]bool) bool {
	for _, winSet := range winningSets {
		win := true
		for action, _ := range winSet {
			_, ok := history[action]
			if !ok {
				win = false
				break
			}
		}
		if win {
			return true
		}
	}
	return false
}

func checkGameResult(history []statemachine.Action) (xwin bool, owin bool) {
	xmoves := map[string]bool{}
	omoves := map[string]bool{}

	if len(history) < 5 {
		return // takes min 5 moves to win
	}

	for i, action := range history {
		if i%2 == 0 { // if X's turn
			xmoves[string(action)[1:]] = true
		} else {
			omoves[string(action)[1:]] = true
		}
	}

	return containsWin(xmoves), containsWin(omoves)
}

func (s StateMachineWithHistory) Clone(state statemachine.StateVector, action string) StateMachineWithHistory {
	m := s.StateMachine.Clone(state)
	sh := StateMachineWithHistory{
		m,
		[]statemachine.Action{},
	}

	for _, a := range s.History {
		sh.History = append(sh.History, a)
	}
	sh.History = append(sh.History, statemachine.Action(action))
	return sh
}

// NOTE: if you do such a test on an unbounded network (like the 'counter' state machine above)
// be warned that golang will happily try to recurse forever
// in the case of a large state space or an unbounded network
// do a random walk simulation to get aggregated result data rather than
// trying to examine the entire state space
func walkNet(sm StateMachineWithHistory) {
	actions, gameOver := sm.ValidActions(1)
	xWin, oWin := checkGameResult(sm.History)

	if !gameOver || xWin || oWin {
		HistoryIndex = append(HistoryIndex, sm.History)
		return
	}

	for action, state := range actions {
		walkNet(sm.Clone(state, action))
	}
}

// the actual test for boundedness is that
// the recursive function walkNet should not infinitely recurse
// this test completes in ~7s
func TestTicTacToeStateSpace(t *testing.T) {
	_, s := statemachine.LoadPnmlFromFile("../examples/octoe.xml")
	// remove extraneous early game-ending actions
	delete(s.Transitions, statemachine.Action("END_O"))
	delete(s.Transitions, statemachine.Action("END_X"))

	t.Run("initialize game", func(t *testing.T) {
		s.Init()
		_, err := s.Commit("EXEC", 1)
		assert.Nil(t, err)
	})

	sm := StateMachineWithHistory{s, []statemachine.Action{}}

	t.Run("walk state space", func(t *testing.T) {
		walkNet(sm)
		games := len(HistoryIndex)

		fmt.Printf("games: %v\n", games)
		assert.Equal(t, games, 255168, "expected count to all possible games")
		testIndex := 77
		fmt.Printf("Game %v: %v\n", testIndex, HistoryIndex[testIndex])
	})

	var XWins []int
	var OWins []int
	var Draws []int

	t.Run("gather stats about moves", func(t *testing.T) {

		for i, history := range HistoryIndex {
			xWin, oWin := checkGameResult(history)

			switch {
			case xWin:
				XWins = append(XWins, i)
			case oWin:
				OWins = append(OWins, i)
			default:
				Draws = append(Draws, i)
			}
		}

		// NOTE: this does not account for symmetry of board states
		assert.Equal(t, 255168, len(XWins)+len(OWins)+len(Draws), "total games should add up")
		assert.Equal(t, 131184, len(XWins), "number of X wins")
		assert.Equal(t, 77904, len(OWins), "number of O wins")
		assert.Equal(t, 46080, len(Draws), "number of Draws")
	})
}
