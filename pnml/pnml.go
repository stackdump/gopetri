package pnml

import (
	"encoding/xml"
	"io/ioutil"
)

type Token struct {
	Id	string  `xml:"id,attr"`
	Blue  int `xml:"blue,attr"`
	Green int `xml:"green,attr"`
	Red   int `xml:"red,attr"`
}
type Name struct {
	Graphics Offset `xml:"graphics>offset"`
	Value    string `xml:"value"`
}

type Capacity struct {
	Value uint64 `xml:"value"`
}

type InitialMarking struct {
	Graphics      Offset `xml:"graphics>position"`
	TokenValueCsv string `xml:"value"`
}

type Coordinates struct {
	X string `xml:"x,attr"`
	Y string `xml:"y,attr"`
}

type Position struct {
	Coordinates `xml:"position"`
}

type Offset struct {
	Coordinates `xml:"offset"`
}

type InfiniteServer struct {
	Value bool `xml:"value"`
}

type Timed struct {
	Value bool `xml:"value"`
}

type Priority struct {
	Value uint64 `xml:"value"`
}

type Orientation struct {
	Value uint64 `xml:"value"`
}

type Rate struct {
	Value float64 `xml:"value"`
}

type ArcPath struct {
	CurvePoint bool `xml:"curvepoint,attr"`
	Id string `xml:"id,attr"`
	X string `xml:"x,attr"`
	Y string `xml:"y,attr"`
}

type ArcInscription struct {
	TokenValueCsv string `xml:"value"`
}

type ArcType struct {
	Value string `xml:"value,attr""`
}

type Arc struct {
	Id 			string	`xml:"id,attr"`
	Source      string	`xml:"source,attr"`
	Target      string 	`xml:"target,attr"`
	ArcPaths	[]ArcPath	`xml:"arcpath"`
	Type		ArcType	 `xml:"type"`
	Inscription  ArcInscription `xml:"inscription"`
}

type Place struct {
	Id             string         `xml:"id,attr"`
	Graphics       Position       `xml:"graphics>position"`
	Name           Name           `xml:"name"`
	Capacity       Capacity       `xml:"capacity"`
	InitialMarking InitialMarking `xml:"initialMarking"`
}

type Transition struct {
	Id             string         `xml:"id,attr"`
	Graphics       Position       `xml:"graphics>position"`
	Name           Name           `xml:"name"`
	InfiniteServer InfiniteServer `xml:"infiniteServer"`
	Timed          Timed          `xml:"timed"`
	Priority       Priority       `xml:"priority"`
	Orientation    Orientation    `xml:"orientation"`
	Rate           Rate           `xml:"rate"`
}

type Net struct {
	Tokens      []Token      `xml:"token"`
	Places      []Place      `xml:"place"`
	Transitions []Transition `xml:"transition"`
	Arcs      	[]Arc      `xml:"arc"`
}


type Pnml struct {
	Nets []Net `xml:"net"`
}

// KLUDGE xml tag needs to be lowercase
type pnml struct {
	Nets []Net `xml:"net"`
}

func LoadFile(path string) (*Pnml, error) {
	p := new(Pnml)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return p, err
	}
	p.Unmarshal(data)
	return p, nil
}

func (p *Pnml) Marshal() ([]byte, error) {
	return xml.Marshal(pnml{p.Nets})
}

func (p *Pnml) Unmarshal(data []byte) error {
	p0 := new(pnml)
	err := xml.Unmarshal(data, p0)
	if err != nil {
		return err
	}

	p.Nets = p0.Nets
	return nil
}
