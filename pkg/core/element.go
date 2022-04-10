package core

import (
	"encoding/json"
	"errors"
	"strings"
)

//EleType is a string representing an element i.e. HYDRO/PYRO/etc...
type EleType int

//ElementType should be pryo, Hydro, Cryo, Electro, Geo, Anemo and maybe dendro
const (
	Electro EleType = iota
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

func (e EleType) String() string {
	return EleTypeString[e]
}

func (e *EleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(EleTypeString[*e])
}

func (e *EleType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range EleTypeString {
		if v == s {
			*e = EleType(i)
			return nil
		}
	}
	return errors.New("unrecognized element")
}

var EleTypeString = [...]string{
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

func StringToEle(s string) EleType {
	for i, v := range EleTypeString {
		if v == s {
			return EleType(i)
		}
	}
	return UnknownElement
}

func EleToDmgP(e EleType) StatType {
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

type ReactionType string

const (
	Overload           ReactionType = "overload"
	Superconduct       ReactionType = "superconduct"
	Freeze             ReactionType = "freeze"
	Melt               ReactionType = "melt"
	Vaporize           ReactionType = "vaporize"
	CrystallizeElectro ReactionType = "crystallize-electro"
	CrystallizeHydro   ReactionType = "crystallize-hydro"
	CrystallizePyro    ReactionType = "crystallize-pyro"
	CrystallizeCryo    ReactionType = "crystallize-cryo"
	SwirlElectro       ReactionType = "swirl-electro"
	SwirlHydro         ReactionType = "swirl-hydro"
	SwirlPyro          ReactionType = "swirl-pyro"
	SwirlCryo          ReactionType = "swirl-cryo"
	ElectroCharged     ReactionType = "electrocharged"
	Shatter            ReactionType = "shatter"
	NoReaction         ReactionType = ""
	FreezeExtend       ReactionType = "FreezeExtend"
)

func absorbCheckWillCollide(p AttackPattern, t Target, index int) bool {
	//shape shouldn't be nil; panic here
	if p.Shape == nil {
		panic("unexpected nil shape")
	}
	//shape can't be nil now, check if type matches
	if !p.Targets[t.Type()] {
		return false
	}

	//check if shape matches
	switch v := p.Shape.(type) {
	case *Circle:
		return t.Shape().IntersectCircle(*v)
	case *Rectangle:
		return t.Shape().IntersectRectangle(*v)
	case *SingleTarget:
		//only true if
		return v.Target == index
	default:
		return false
	}
}

func (c *Core) AbsorbCheck(p AttackPattern, prio ...EleType) EleType {

	// check targets for collision first

	for _, e := range prio {
		for i, t := range c.Targets {
			if absorbCheckWillCollide(p, t, i) && t.AuraContains(e) {
				c.Log.NewEvent(
					"infusion check picked up "+e.String(),
					LogElementEvent,
					-1,
				)
				return e
			}
		}
	}
	return NoElement
}
