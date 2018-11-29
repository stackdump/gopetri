package pnml_test

import (
	. "github.com/stackdump/gopetri/pnml"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	p, err := LoadFile("../examples/counter.xml")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Nets) != 1 {
		t.Fatal("failed to load xml file")
	}
	data, _ := p.Marshal()
	println(string(data))
}