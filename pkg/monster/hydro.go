package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

type AuraHydro struct {
	*Element
}

func (a *AuraHydro) React(ds *def.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case def.Anemo:
		//this one doesn't do aoe damage for some reason
		ds.ReactionType = def.SwirlHydro
		//queue swirl dmg
		t.queueReaction(ds, def.SwirlHydro, a.CurrentDurability, 1)
		//reduce hydro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case def.Geo:
		ds.ReactionType = def.CrystallizeHydro
		//crystallize adds shield
		shd := NewCrystallizeShield(def.Hydro, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case def.Pyro:
		//vaporize, pyro into hydro = weak
		ds.ReactionType = def.Vaporize
		ds.ReactMult = 1.5
		a.Reduce(ds, 0.5)
	case def.Hydro:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		return a, false
	case def.Cryo:
		//first reduce hydro durability by incoming cryo; capped at existing
		red := a.Reduce(ds, 1)
		if a.CurrentDurability < 0 {
			a = nil
		}
		ds.ReactionType = def.Freeze
		//since cryo is applied, cryo aura is nil
		return newFreeze(nil, a, red, t, ds, t.sim.Frame()), true
	case def.Electro:
		//ec
		e := &AuraElectro{}
		e.Element = &Element{}
		e.T = def.Electro
		e.Attach(ds.Durability, t.sim.Frame())
		ds.ReactionType = def.ElectroCharged
		return newEC(e, a, t, ds, t.sim.Frame()), true
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, true
	}
	return a, true
}
