package monster

import "github.com/genshinsim/gsim/pkg/def"

type AuraElectro struct {
	*Element
}

func (a *AuraElectro) React(ds *def.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case def.Anemo:
		ds.ReactionType = def.SwirlElectro
		//queue swirl dmg
		t.queueReaction(ds, def.SwirlElectro, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case def.Geo:
		ds.ReactionType = def.CrystallizeElectro
		//crystallize adds shield
		shd := NewCrystallizeShield(def.Electro, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case def.Pyro:
		//overload + reduce
		ds.ReactionType = def.Overload
		t.queueReaction(ds, def.Overload, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	case def.Hydro:
		//ec
		//don't touch they hydro aura, instead, queue up ec ticks
		h := &AuraHydro{}
		h.Element = &Element{}
		h.T = def.Hydro
		h.Attach(ds.Durability, t.sim.Frame())
		ds.ReactionType = def.ElectroCharged
		return newEC(a, h, t, ds, t.sim.Frame()), true
	case def.Cryo:
		//superconduct
		ds.ReactionType = def.Superconduct
		t.queueReaction(ds, def.Superconduct, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	case def.Electro:
		//have to be careful here.. if there's also a hydro then we have to trigger a tick
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, true
	}
	return a, true
}

/**
t=0; + electro(25) = [`electro(20)`] //initial electro application, t in frames
t=2; [`electro(20)`] + hydro(25) = [`electro(20)`, `hydro(20)`] //hydro
t=3; ec damage incurred, gauge reduction queued for t = 3 + 6 (0.1 seconds)
t=5; [`electro(19)`, `hydro(19)`] + hydro(25) = [`electro(19)`, `hydro(20)`] //second hydro, assume minor decay
t=5; ec damage incurred, gauge reduction queued for t = 5 + 6 (0.1 seconds)
t=9; [`electro(19)`, `hydro(20)`] -> gauge reduction -> [`electro(9)`, `hydro(10)`]
t=11; [`electro(9)`, `hydro(10)`] -> gauge reduction -> all gone
**/
