package monster

import "github.com/genshinsim/gcsim/pkg/core"

type AuraElectro struct {
	*Element
}

func (a *AuraElectro) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}

	reactionTriggered := true
	switch ds.Element {
	case core.Anemo:
		ds.ReactionType = core.SwirlElectro
		//queue swirl dmg
		t.queueReaction(ds, core.SwirlElectro, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case core.Geo:
		ds.ReactionType = core.CrystallizeElectro
		//crystallize adds shield
		shd := NewCrystallizeShield(core.Electro, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case core.Pyro:
		//overload + reduce
		ds.ReactionType = core.Overload
		t.queueReaction(ds, core.Overload, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	case core.Hydro:
		//ec
		//don't touch they hydro aura, instead, queue up ec ticks
		h := &AuraHydro{}
		h.Element = &Element{}
		h.T = core.Hydro
		h.Attach(ds.Durability, t.core.F)
		ds.ReactionType = core.ElectroCharged
		return newEC(a, h, t, ds, t.core.F), true
	case core.Cryo:
		//superconduct
		ds.ReactionType = core.Superconduct
		t.queueReaction(ds, core.Superconduct, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	case core.Electro:
		//have to be careful here.. if there's also a hydro then we have to trigger a tick
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		reactionTriggered = false
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, reactionTriggered
	}
	return a, reactionTriggered
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
