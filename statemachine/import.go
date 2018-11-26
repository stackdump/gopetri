package statemachine

import (
	"github.com/stackdump/gopetri/pnml"
)

func LoadPnmlFromFile(path string) StateMachine {
	net, err := pnml.LoadFile(path)
	_ = net
	_ = err
	// fixme
	return StateMachine{}
}
