package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

type AuraFrozen struct {
	*Element

	hydro    *AuraHydro
	cryo     *AuraCyro
	source   int
	snapshot def.Snapshot
	t        *Target
}

func newFreeze(c *AuraCyro, h *AuraHydro, dur def.Durability, t *Target, ds *def.Snapshot, f int) Aura {
	fz := AuraFrozen{}
	fz.Element = &Element{}
	fz.T = def.Frozen
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

func (a *AuraFrozen) React(ds *def.Snapshot, t *Target) (Aura, bool) {
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
		if ds.Durability > 0 && a.hydro != nil {
			next, _ := a.hydro.React(ds, t)
			a.hydro = next.(*AuraHydro)
		}
	case def.Geo:
		ds.ReactionType = def.CrystallizeCryo
		//crystallize adds shield
		shd := NewCrystallizeShield(def.Cryo, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
		//not sure if we can proc second crystallize?
	case def.Pyro:
		//melt, pyro into cryo = strong; should be melt only
		ds.ReactionType = def.Melt
		ds.ReactMult = 2
		a.Reduce(ds, 2)
	case def.Hydro:
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
			ds.ReactionType = def.Freeze
			return a, true
		}
		//check if we're topping up hydro here
		if a.hydro != nil {
			a.hydro.Refresh(ds.Durability)
			ds.Durability = 0
			return a, false
		}
		//otherwise attach it
		r := &AuraHydro{}
		r.Element = &Element{}
		r.T = def.Hydro
		r.Attach(ds.Durability, t.sim.Frame())
		a.hydro = r
		return a, false
	case def.Cryo:
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
			ds.ReactionType = def.Freeze
			return a, true
		}
		//check if we're topping up hydro here
		if a.cryo != nil {
			a.cryo.Refresh(ds.Durability)
			ds.Durability = 0
			return a, false
		}
		//otherwise attach it
		r := &AuraCyro{}
		r.Element = &Element{}
		r.T = def.Cryo
		r.Attach(ds.Durability, t.sim.Frame())
		a.cryo = r
		return a, false
	case def.Electro:
		//superconduct
		ds.ReactionType = def.Superconduct
		t.queueReaction(ds, def.Superconduct, a.CurrentDurability, 1)
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
