package monster

import "github.com/genshinsim/gsim/pkg/def"

type PyroAura struct {
	*Element
}

func (p *PyroAura) React(ds *def.Snapshot, s def.Sim, auras []Aura) bool {
	if ds.Durability == 0 {
		return false
	}

	switch ds.Element {
	case def.Anemo:
		//swirl dmg
	case def.Geo:
		//crystallize dmg
	case def.Pyro:
		//refresh
		p.Element.Refresh(ds.Durability)
		ds.Durability = 0
		return true
	case def.Hydro:
		//vaporize + reduce
	case def.Cryo:
		//melt + reduce
	case def.Electro:
		//overload + reduce

	default:
		return false
	}
	return true
}
