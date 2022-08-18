package attributes

import (
	"encoding/json"
	"errors"
	"strings"
)

//Element is a string representing an element i.e. HYDRO/PYRO/etc...
type Element int

//ElementType should be pryo, Hydro, Cryo, Electro, Geo, Anemo and maybe dendro
const (
	Electro Element = iota
	Pyro
	Cryo
	Hydro
	Frozen
	ElementDelimAttachable
	Anemo
	Dendro
	Geo
	NoElement
	Physical
	EC
	UnknownElement
	EndEleType
)

func (e Element) String() string {
	return ElementString[e]
}

func (e Element) MarshalJSON() ([]byte, error) {
	return json.Marshal(ElementString[e])
}

func (e *Element) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range ElementString {
		if v == s {
			*e = Element(i)
			return nil
		}
	}
	return errors.New("unrecognized element")
}

var ElementString = [...]string{
	"electro",
	"pyro",
	"cryo",
	"hydro",
	"frozen",
	"delim",
	"anemo",
	"dendro",
	"geo",
	"",
	"physical",
	"ec",
	"unknown",
}

func StringToEle(s string) Element {
	for i, v := range ElementString {
		if v == s {
			return Element(i)
		}
	}
	return UnknownElement
}

func EleToDmgP(e Element) Stat {
	switch e {
	case Anemo:
		return AnemoP
	case Cryo:
		return CryoP
	case Electro:
		return ElectroP
	case Geo:
		return GeoP
	case Hydro:
		return HydroP
	case Pyro:
		return PyroP
	case Dendro:
		return DendroP
	case Physical:
		return PhyP
	}
	return -1
}
