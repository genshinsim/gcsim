package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryFreeze(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	switch a.Info.Element {
	case attributes.Hydro:
		//if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[attributes.Cryo] > ZeroDur {
			consumed := r.triggerFreeze(r.Durability[attributes.Cryo], a.Info.Durability)
			r.Durability[attributes.Cryo] -= consumed
			r.Durability[attributes.Cryo] = max(r.Durability[attributes.Cryo], 0)
			//TODO: we're not setting src durability to zero here but should be ok b/c no reaction comes after freeze
			//ec should have been taken care of already
			a.Info.Durability -= consumed
			a.Info.Durability = max(a.Info.Durability, 0)
			r.core.Events.Emit(event.OnFrozen, r.self, a)
			return
		}
		//otherwise attach hydro only if frozen exists
		if r.Durability[attributes.Frozen] < ZeroDur {
			return
		}
		//try refill first - this will use up all durability if ok
		r.tryRefill(attributes.Hydro, &a.Info.Durability)
		//otherwise attach
		r.tryAttach(attributes.Hydro, &a.Info.Durability)
	case attributes.Cryo:
		if r.Durability[attributes.Hydro] > ZeroDur {
			consumed := r.triggerFreeze(r.Durability[attributes.Hydro], a.Info.Durability)
			r.Durability[attributes.Hydro] -= consumed
			r.Durability[attributes.Hydro] = max(r.Durability[attributes.Hydro], 0)
			a.Info.Durability -= consumed
			a.Info.Durability = max(a.Info.Durability, 0)
			r.core.Events.Emit(event.OnFrozen, r.self, a)
			return
		}
		//otherwise attach cryo only if frozen exists
		if r.Durability[attributes.Frozen] < ZeroDur {
			return
		}
		//try refill first - this will use up all durability if ok
		r.tryRefill(attributes.Cryo, &a.Info.Durability)
		//otherwise attach
		r.tryAttach(attributes.Cryo, &a.Info.Durability)
	default:
		//should be here
		return
	}
}

func max(a, b combat.Durability) combat.Durability {
	if a > b {
		return a
	}
	return b
}

func min(a, b combat.Durability) combat.Durability {
	if a > b {
		return b
	}
	return a
}

func (r *Reactable) ShatterCheck(a *combat.AttackEvent) {
	if a.Info.StrikeType != combat.StrikeTypeBlunt || r.Durability[attributes.Frozen] < ZeroDur {
		return
	}
	//remove 200 freeze gauge if availabe
	r.Durability[attributes.Frozen] -= 200
	r.checkFreeze()
	//trigger shatter attack
	ai := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(combat.Shatter),
		AttackTag:        combat.AttackTagShatter,
		ICDTag:           combat.ICDTagShatter,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          attributes.Physical,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	ai.FlatDmg = 1.5 * r.calcReactionDmg(ai, em)
	//shatter is a self attack
	r.core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(r.self.Index(), r.self.Type()),
		-1,
		1,
	)

}

// add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b combat.Durability) combat.Durability {
	d := min(a, b)
	//trigger freeze should only addDurability and should not touch decay rate
	r.addDurability(attributes.Frozen, 2*d)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[attributes.Frozen] <= ZeroDur {
		r.Durability[attributes.Frozen] = 0
		r.core.Events.Emit(event.OnAuraDurabilityDepleted, r.self, attributes.Frozen)
		//trigger another attack here, purely for the purpose of breaking bubbles >.>
		ai := combat.AttackInfo{
			ActorIndex:  0,
			DamageSrc:   r.self.Key(),
			Abil:        "Freeze Broken",
			AttackTag:   combat.AttackTagNone,
			ICDTag:      combat.ICDTagNone,
			ICDGroup:    combat.ICDGroupDefault,
			Element:     attributes.NoElement,
			SourceIsSim: true,
			DoNotLog:    true,
		}
		//TODO: delay attack by 1 frame ok?
		r.core.QueueAttack(ai, combat.NewDefSingleTarget(r.self.Index(), r.self.Type()), -1, 1)
	}
}
