package statemachine

import (
	"github.com/stackdump/gopetri/pnml"
	"strconv"
	"strings"
)

type Place struct {
	Initial  uint64
	Offset   int
	Capacity uint64
}

type PetriNet struct {
	Places      map[string]Place
	Transitions map[Action]Transition
	Pnml        *pnml.Pnml
}

func (p PetriNet) getOffset(placeName string) (int, bool) {
	for placeID, place := range p.Places {
		if placeID == placeName {
			return place.Offset, true
		}
	}
	return -1, false
}

func getWeight(a pnml.Arc) int64 {
	tokenVals := strings.Split(a.Inscription.TokenValueCsv, ",")
	val, err := strconv.ParseInt(tokenVals[1], 10, 64)

	if err != nil || tokenVals[0] != "Default" {
		panic("Error Parsing Token Weight")
	}
	return val
}

func GetInitialValue(m pnml.InitialMarking) uint64 {
	tokenVals := strings.Split(m.TokenValueCsv, ",")
	val, err := strconv.ParseInt(tokenVals[1], 10, 64)

	if err != nil || tokenVals[0] != "Default" {
		panic("Error Parsing InitialMarking")
	}
	return uint64(val)
}

func (p PetriNet) GetEmptyVector() []int64 {
	emptyVector := []int64{}

	for x := 0; x < len(p.Places); x++ {
		emptyVector = append(emptyVector, int64(0))
	}
	return emptyVector
}

func (p PetriNet) GetEmptyState() []uint64 {
	emptyState := []uint64{}

	for x := 0; x < len(p.Places); x++ {
		emptyState = append(emptyState, uint64(0))
	}
	return emptyState
}

func (p PetriNet) GetInitialState() StateVector {
	if p.Places == nil || len(p.Places) == 0 {
		panic("no places defined")
	}
	init := p.GetEmptyState()
	for _, place := range p.Places {
		init[place.Offset] = place.Initial
	}
	return StateVector(init[:])
}

func (p PetriNet) GetCapacityVector() StateVector {
	cap := p.GetEmptyState()
	for _, place := range p.Places {
		cap[place.Offset] = place.Capacity
	}
	return StateVector(cap[:])
}

func (p PetriNet) StateMachine() StateMachine {
	return StateMachine{
		Initial:     p.GetInitialState(),
		Capacity:    p.GetCapacityVector(),
		Transitions: p.Transitions,
		State:       StateVector{},
	}
}

func LoadPnmlFromFile(path string) (PetriNet, StateMachine) {
	pnmlDef, _ := pnml.LoadFile(path)

	petriNet := PetriNet{
		map[string]Place{},
		map[Action]Transition{},
		pnmlDef,
	}

	if len(pnmlDef.Nets) != 1 {
		panic("Expect a single petri-net definition per file")
	}

	net := pnmlDef.Nets[0]

	for offset, p := range net.Places {
		petriNet.Places[p.Id] =
			Place{
				Initial:  GetInitialValue(p.InitialMarking),
				Offset:   offset,
				Capacity: p.Capacity.Value,
			}
	}

	for _, txn := range net.Transitions {
		petriNet.Transitions[Action(txn.Id)] = petriNet.GetEmptyVector()
	}

	for _, arc := range net.Arcs {
		var action string
		var sign int64
		var offset int

		targetOffset, targetIsPlace := petriNet.getOffset(arc.Target)
		sourceOffset, sourceIsPlace := petriNet.getOffset(arc.Source)

		if sourceIsPlace {
			action = arc.Target
			offset = sourceOffset
			sign = -1
		}

		if targetIsPlace {
			action = arc.Source
			offset = targetOffset
			sign = 1
		}

		petriNet.Transitions[Action(action)][offset] += sign * getWeight(arc)
	}

	return petriNet, petriNet.StateMachine()
}
