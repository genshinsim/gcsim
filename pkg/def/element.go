package def

//EleType is a string representing an element i.e. HYDRO/PYRO/etc...
type EleType int

//ElementType should be pryo, Hydro, Cryo, Electro, Geo, Anemo and maybe dendro
const (
	Electro EleType = iota
	Pyro
	Anemo
	Cryo
	Frozen
	Hydro
	Dendro
	Geo
	NoElement
	ElementMaxCount
	Physical
	UnknownElement
)

func (e EleType) String() string {
	return EleTypeString[e]
}

var EleTypeString = [...]string{
	"electro",
	"pyro",
	"anemo",
	"cryo",
	"frozen",
	"hydro",
	"dendro",
	"geo",
	"",
	"delim",
	"physical",
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
