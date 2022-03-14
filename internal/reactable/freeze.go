package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (r *Reactable) tryFreeze(a *coretype.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	switch a.Info.Element {
	case core.Hydro:
		//if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[coretype.Cryo] > ZeroDur {
			consumed := r.triggerFreeze(r.Durability[coretype.Cryo], a.Info.Durability)
			r.Durability[coretype.Cryo] -= consumed
			r.Durability[coretype.Cryo] = max(r.Durability[coretype.Cryo], 0)
			//TODO: we're not setting src durability to zero here but should be ok b/c no reaction comes after freeze
			//ec should have been taken care of already
			a.Info.Durability -= consumed
			a.Info.Durability = max(a.Info.Durability, 0)
			return
		}
		//otherwise attach hydro only if frozen exists
		if r.Durability[coretype.Frozen] < ZeroDur {
			return
		}
		//try refill first - this will use up all durability if ok
		r.tryRefill(core.Hydro, &a.Info.Durability)
		//otherwise attach
		r.tryAttach(core.Hydro, &a.Info.Durability)
	case coretype.Cryo:
		if r.Durability[core.Hydro] > ZeroDur {
			consumed := r.triggerFreeze(r.Durability[core.Hydro], a.Info.Durability)
			r.Durability[core.Hydro] -= consumed
			r.Durability[core.Hydro] = max(r.Durability[core.Hydro], 0)
			a.Info.Durability -= consumed
			a.Info.Durability = max(a.Info.Durability, 0)
			return
		}
		//otherwise attach cryo only if frozen exists
		if r.Durability[coretype.Frozen] < ZeroDur {
			return
		}
		//try refill first - this will use up all durability if ok
		r.tryRefill(coretype.Cryo, &a.Info.Durability)
		//otherwise attach
		r.tryAttach(coretype.Cryo, &a.Info.Durability)
	default:
		//should be here
		return
	}
	r.core.Events.Emit(core.OnFrozen, r.self, a)
}

func max(a, b core.Durability) core.Durability {
	if a > b {
		return a
	}
	return b
}

func min(a, b core.Durability) core.Durability {
	if a > b {
		return b
	}
	return a
}

func (r *Reactable) ShatterCheck(a *coretype.AttackEvent) {
	if a.Info.StrikeType != core.StrikeTypeBlunt || r.Durability[coretype.Frozen] < ZeroDur {
		return
	}
	//remove 200 freeze gauge if availabe
	r.Durability[coretype.Frozen] -= 200
	r.checkFreeze()
	//trigger shatter attack
	ai := core.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(core.Shatter),
		AttackTag:        core.AttackTagShatter,
		ICDTag:           core.ICDTagShatter,
		ICDGroup:         core.ICDGroupReactionA,
		Element:          core.Physical,
		IgnoreDefPercent: 1,
	}
	ai.FlatDmg = 1.5 * r.calcReactionDmg(ai)
	//shatter is a self attack
	r.core.Combat.QueueAttack(
		ai,
		core.NewDefSingleTarget(r.self.Index(), r.self.Type()),
		-1,
		1,
	)

}

//add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b core.Durability) core.Durability {
	d := min(a, b)
	//trigger freeze should only addDurability and should not touch decay rate
	r.addDurability(coretype.Frozen, 2*d)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[coretype.Frozen] <= ZeroDur {
		r.Durability[coretype.Frozen] = 0
		r.core.Events.Emit(core.OnAuraDurabilityDepleted, r.self, coretype.Frozen)
		//trigger another attack here, purely for the purpose of breaking bubbles >.>
		ai := core.AttackInfo{
			ActorIndex:  0,
			DamageSrc:   r.self.Index(),
			Abil:        "Freeze Broken",
			AttackTag:   core.AttackTagNone,
			ICDTag:      core.ICDTagNone,
			ICDGroup:    core.ICDGroupDefault,
			Element:     core.NoElement,
			SourceIsSim: true,
			DoNotLog:    true,
		}
		//TODO: delay attack by 1 frame ok?
		r.core.Combat.QueueAttack(ai, core.NewDefSingleTarget(r.self.Index(), r.self.Type()), -1, 1)
	}
}
