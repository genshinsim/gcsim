package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type AuraPyro struct {
	*Element
}

func (a *AuraPyro) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}

	reactionTriggered := true
	switch ds.Element {
	case core.Anemo:
		ds.ReactionType = core.SwirlPyro
		//queue swirl dmg
		t.queueReaction(ds, core.SwirlPyro, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case core.Geo:
		ds.ReactionType = core.CrystallizePyro
		//crystallize adds shield
		shd := NewCrystallizeShield(core.Pyro, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case core.Pyro:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		reactionTriggered = false
	case core.Hydro:
		//vaporize + reduce
		ds.ReactionType = core.Vaporize
		ds.ReactMult = 2
		ds.IsMeltVape = true
		a.Reduce(ds, 2)
	case core.Cryo:
		//melt + reduce
		ds.ReactionType = core.Melt
		ds.ReactMult = 1.5
		ds.IsMeltVape = true
		//vaporize + reduce
		a.Reduce(ds, 0.5)
	case core.Electro:
		//overload + reduce
		ds.ReactionType = core.Overload
		t.queueReaction(ds, core.Overload, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, reactionTriggered
	}
	return a, reactionTriggered
}
