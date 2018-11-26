package statemachine_test

import (
	"testing"
	"github.com/stackdump/gopetri/statemachine"
)

func TestCounterMachine(t *testing.T) {
	s := statemachine.LoadPnmlFromFile("../examples/counter.xml")
	if len(s.Transitions) == 0 {
		t.Fatal("failed to load statemachine")
	}
}

