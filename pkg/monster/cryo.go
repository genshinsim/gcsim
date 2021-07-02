package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

type AuraCyro struct {
	*Element
}

func (a *AuraCyro) React(ds *def.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case def.Anemo:
		ds.ReactionType = def.SwirlCryo
		//queue swirl dmg
		t.queueReaction(ds, def.SwirlCryo, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case def.Geo:
		ds.ReactionType = def.CrystallizeCryo
		//crystallize adds shield
		shd := NewCrystallizeShield(def.Cryo, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case def.Pyro:
		//melt, pyro into cryo = strong
		ds.ReactionType = def.Melt
		ds.ReactMult = 2
		a.Reduce(ds, 2)
	case def.Hydro:
		//first reduce hydro durability by incoming cryo; capped at existing
		red := a.Reduce(ds, 1)
		if a.CurrentDurability < 0 {
			a = nil
		}
		ds.ReactionType = def.Freeze
		//since cryo is applied, cryo aura is nil
		return newFreeze(a, nil, red, t, ds, t.sim.Frame()), true
	case def.Cryo:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
	case def.Electro:
		//superconduct
		ds.ReactionType = def.Superconduct
		t.queueReaction(ds, def.Superconduct, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, true
	}
	return a, true
}
