package statemachine

import (
	"github.com/stackdump/gopetri/pnml"
)

type Place struct {
	Initial uint64
	Offset int
}

type Txn struct {
	Delta []int64
}

type PetriNet struct {
	Places map[string]Place
	Transitions map[string]Txn
}

// TODO add methods to convert to a vector state machine
func LoadPnmlFromFile(path string) StateMachine {
	petriNet := PetriNet{
		map[string]Place{},
		map[string]Txn{},
	}
	pnmlDef, _ := pnml.LoadFile(path)
	net := pnmlDef.Nets[0]
	var emptyVector []int64
	for x := 0; x < len(net.Places); x++ {
		emptyVector = append(emptyVector, int64(0))
	}

	for offset, p := range net.Places {
		petriNet.Places[p.Id] =
			Place{
				Initial: 0,
				Offset: offset,
			}
	}

	for _, txn := range net.Transitions {
		petriNet.Transitions[txn.Id] =
			Txn{
				Delta: emptyVector,
			}
	}
	_ = net
	return StateMachine{}
}
