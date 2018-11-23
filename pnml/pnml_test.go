package pnml_test

import (
	"github.com/go-xmlfmt/xmlfmt"
	. "github.com/stackdump/gopetri/pnml"
	"testing"
)

var net0 Net = Net{
	[]Token{
		Token{
			Blue:  0,
			Green: 0,
			Red:   0,
		},
	},
	[]Place{
		Place{
			Name: Name{
				Offset{},
				"foo",
			},
			Capacity: Capacity{1},
			InitialMarking: InitialMarking{
				Offset{Coordinates{"0.0", "0.0"}},
				"Default,0",
			},
		},
	},
	[]Transition{
		Transition{
			Id:       "t0",
			Graphics: Position{Coordinates{"0.0", "0.0"}},
			Name: Name{
				Offset{
					Coordinates{"0.0", "0.0"},
				},
				"bar",
			},
			InfiniteServer: InfiniteServer{Value: false},
			Timed:          Timed{Value: false},
			Priority:       Priority{Value: 1},
			Orientation:    Orientation{Value: 0},
			Rate:           Rate{Value: 1.0},
		},
	},
}

// print out formatted XML
func TestMarsharToXML(t *testing.T) {
	var data []byte
	var err error

	t.Run("marshal and unmarshal", func(t *testing.T) {
		p1 := NewPnml([]Net{net0})
		data, err = p1.Marshal()
		p2 := new(Pnml)
		err = p2.Unmarshal(data)

		if len(p2.Nets) != 1 {
			t.Fatal("failed to unmarshal xml")
		}
	})

	// just print XML
	x := xmlfmt.FormatXML(string(data), "\t", "  ")
	if err != nil {
		t.Fatal(err)
	}
	println(x)
}
