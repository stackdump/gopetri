package statemachine_test

import (
	"fmt"
	"testing"
	"github.com/stackdump/gopetri/statemachine"
)

func TestCounterMachine(t *testing.T) {
	s := statemachine.LoadPnmlFromFile("../examples/counter.xml")

	t.Run("load machine ", func(t *testing.T ) {
		s.Init()
		if len(s.Transitions) == 0 {
			t.Fatal("failed to load statemachine")
		}
		fmt.Printf("%v\n", s)
	})

	t.Run("commit actions", func(t *testing.T ){
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

	t.Run("check bounds", func(t *testing.T ) {
		// KLUDGE: setting Capacity to all 0 will not restrict transactions
		// during execution but will cause the manual check to fail
		if s.InBounds() {
			t.Fatalf("expeced to be out-of-bounds")
		}
		fmt.Printf("%v\n", s)
	})

}
