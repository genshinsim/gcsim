package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type AuraHydro struct {
	*Element
}

func (a *AuraHydro) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}

	reactionTriggered := true
	switch ds.Element {
	case core.Anemo:
		//this one doesn't do aoe damage for some reason
		ds.ReactionType = core.SwirlHydro
		//queue swirl dmg
		t.queueReaction(ds, core.SwirlHydro, a.CurrentDurability, 1)
		//reduce hydro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case core.Geo:
		ds.ReactionType = core.CrystallizeHydro
		//crystallize adds shield
		shd := NewCrystallizeShield(core.Hydro, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case core.Pyro:
		//vaporize, pyro into hydro = weak
		ds.ReactionType = core.Vaporize
		ds.ReactMult = 1.5
		ds.IsMeltVape = true
		a.Reduce(ds, 0.5)
	case core.Hydro:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		reactionTriggered = false
	case core.Cryo:
		//first reduce hydro durability by incoming cryo; capped at existing
		red := a.Reduce(ds, 1)
		if a.CurrentDurability < 0 {
			a = nil
		}
		ds.ReactionType = core.Freeze
		//since cryo is applied, cryo aura is nil
		return newFreeze(nil, a, red, t, ds, t.core.F), true
	case core.Electro:
		//ec
		e := &AuraElectro{}
		e.Element = &Element{}
		e.T = core.Electro
		e.Attach(ds.Durability, t.core.F)
		ds.ReactionType = core.ElectroCharged
		return newEC(e, a, t, ds, t.core.F), true
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, reactionTriggered
	}
	return a, reactionTriggered
}
