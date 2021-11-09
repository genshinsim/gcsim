package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type AuraFrozen struct {
	*Element

	hydro    *AuraHydro
	cryo     *AuraCyro
	source   int
	snapshot core.Snapshot
	t        *Target
}

func (a *AuraFrozen) AuraContains(ele ...core.EleType) bool {
	for _, v := range ele {
		if v == core.Hydro && a.hydro != nil {
			return true
		}
		if v == core.Cryo && a.cryo != nil {
			return true
		}
		if v == core.Frozen {
			return true
		}
	}
	return false
}

func newFreeze(c *AuraCyro, h *AuraHydro, dur core.Durability, t *Target, ds *core.Snapshot, f int) Aura {
	fz := AuraFrozen{}
	fz.Element = &Element{}
	fz.T = core.Frozen
	//one of these 2 should be nil
	fz.cryo = c
	fz.hydro = h
	fz.source = f
	fz.snapshot = ds.Clone()
	fz.t = t

	fz.Start = f
	fz.MaxDurability = 2 * dur
	fz.CurrentDurability = 2 * dur

	//durability formula: durability = -1.25t^2 + (d - 10t) <- t is in seconds
	//in frames we get -1.25(t/60)^2 + (d - 10(t/60))
	//dx/dy = -2.50(t/60) - 10/60
	//dx2/dy = -2.50/60
	//at t = 0, decay = 0?
	fz.DecayRate = 1.0 / 6.0

	return &fz
}

func (a *AuraFrozen) Tick() bool {
	//tick down hydro or cryo if not nil
	if a.hydro != nil {
		done := a.hydro.Tick()
		if done {
			a.hydro = nil
		}
	}
	if a.cryo != nil {
		done := a.cryo.Tick()
		if done {
			a.cryo = nil
		}
	}
	//reduce durability by decay, then increase decay by (1.25/60)
	a.DecayRate += 1.0 / 1440.0
	// log.Printf("current dur %v; decay rate %v\n", a.CurrentDurability, a.DecayRate)
	a.CurrentDurability -= a.DecayRate

	if a.CurrentDurability <= 0 {
		//check if hydro or cryo is there, if there set that to target aura, other wise return true
		if a.cryo != nil {
			a.t.aura = a.cryo
			return false
		}
		if a.hydro != nil {
			a.t.aura = a.hydro
			return false
		}
		return true
	}

	return false
}

func (a *AuraFrozen) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case core.Anemo:
		if a.hydro != nil {
			next, _ := a.hydro.React(ds, t)
			a.hydro = next.(*AuraHydro)
			ds.ReactionType = core.SwirlHydro
		}
		if ds.Durability > 0 {
			ds.ReactionType = core.SwirlCryo
			//queue swirl dmg
			t.queueReaction(ds, core.SwirlCryo, a.CurrentDurability, 1)
			//reduce pyro by 0.5 of anemo
			a.Reduce(ds, 0.5)
		}
	case core.Geo:
		ds.ReactionType = core.CrystallizeCryo
		//crystallize adds shield
		shd := NewCrystallizeShield(core.Cryo, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
		//not sure if we can proc second crystallize?
	case core.Pyro:
		//melt, pyro into cryo = strong; should be melt only
		ds.ReactionType = core.Melt
		ds.ReactMult = 2
		ds.IsMeltVape = true
		a.Reduce(ds, 2)
	case core.Hydro:
		//check if we need to refresh freeze
		if a.cryo != nil {
			//top up to the original max durability
			red := a.cryo.Reduce(ds, 1)
			if a.cryo.CurrentDurability < 0 {
				a.cryo = nil
			}
			a.CurrentDurability += red
			if a.CurrentDurability > a.MaxDurability {
				a.CurrentDurability = a.MaxDurability
			}
			ds.ReactionType = core.Freeze
			return a, true
		}
		//check if we're topping up hydro here
		if a.hydro != nil {
			a.hydro.Refresh(ds.Durability)
			ds.Durability = 0
			return a, true
		}
		//otherwise attach it
		r := &AuraHydro{}
		r.Element = &Element{}
		r.T = core.Hydro
		r.Attach(ds.Durability, t.core.F)
		a.hydro = r
	case core.Cryo:
		//check if we need to refresh freeze
		if a.hydro != nil {
			//top up to the original max durability
			red := a.hydro.Reduce(ds, 1)
			if a.hydro.CurrentDurability < 0 {
				a.hydro = nil
			}
			a.CurrentDurability += red
			if a.CurrentDurability > a.MaxDurability {
				a.CurrentDurability = a.MaxDurability
			}
			ds.ReactionType = core.Freeze
			return a, true
		}
		//check if we're topping up hydro here
		if a.cryo != nil {
			a.cryo.Refresh(ds.Durability)
			ds.Durability = 0
			return a, true
		}
		//otherwise attach it
		r := &AuraCyro{}
		r.Element = &Element{}
		r.T = core.Cryo
		r.Attach(ds.Durability, t.core.F)
		a.cryo = r
	case core.Electro:
		//superconduct
		ds.ReactionType = core.Superconduct
		t.queueReaction(ds, core.Superconduct, a.CurrentDurability, 1)
		a.Reduce(ds, 1)
		if ds.Durability > 0 && a.hydro != nil {
			next, _ := a.hydro.React(ds, t)
			a.hydro = next.(*AuraHydro)
		}
	default:
		return a, false
	}
	if a.CurrentDurability < 0 {
		return nil, true
	}
	return a, true
}
