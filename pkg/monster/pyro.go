package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

type AuraPyro struct {
	*Element
}

func (a *AuraPyro) React(ds *def.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}

	switch ds.Element {
	case def.Anemo:
		ds.ReactionType = def.SwirlPyro
		//queue swirl dmg
		t.queueReaction(ds, def.SwirlPyro, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case def.Geo:
		ds.ReactionType = def.CrystallizePyro
		//crystallize adds shield
		shd := NewCrystallizeShield(def.Pyro, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case def.Pyro:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		return a, false
	case def.Hydro:
		//vaporize + reduce
		ds.ReactionType = def.Vaporize
		ds.ReactMult = 2
		a.Reduce(ds, 2)
	case def.Cryo:
		//melt + reduce
		ds.ReactionType = def.Melt
		ds.ReactMult = 1.5
		//vaporize + reduce
		a.Reduce(ds, 0.5)
	case def.Electro:
		//overload + reduce
		ds.ReactionType = def.Overload
		t.queueReaction(ds, def.Overload, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, true
	}
	return a, true
}
