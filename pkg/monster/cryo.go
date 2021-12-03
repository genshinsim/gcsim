package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type AuraCyro struct {
	*Element
}

func (a *AuraCyro) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}

	reactionTriggered := true
	switch ds.Element {
	case core.Anemo:
		ds.ReactionType = core.SwirlCryo
		//queue swirl dmg
		t.queueReaction(ds, core.SwirlCryo, a.CurrentDurability, 1)
		//reduce pyro by 0.5 of anemo
		a.Reduce(ds, 0.5)
	case core.Geo:
		ds.ReactionType = core.CrystallizeCryo
		//crystallize adds shield
		shd := NewCrystallizeShield(core.Cryo, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case core.Pyro:
		//melt, pyro into cryo = strong
		ds.ReactionType = core.Melt
		ds.ReactMult = 2
		ds.IsMeltVape = true
		a.Reduce(ds, 2)
	case core.Hydro:
		//first reduce hydro durability by incoming cryo; capped at existing
		red := a.Reduce(ds, 1)
		if a.CurrentDurability < 0 {
			a = nil
		}
		ds.ReactionType = core.Freeze
		//since cryo is applied, cryo aura is nil
		return newFreeze(a, nil, red, t, ds, t.core.F), true
	case core.Cryo:
		//refresh
		a.Refresh(ds.Durability)
		ds.Durability = 0
		reactionTriggered = false
	case core.Electro:
		//superconduct
		ds.ReactionType = core.Superconduct
		t.queueReaction(ds, core.Superconduct, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, reactionTriggered
	}
	return a, reactionTriggered
}
