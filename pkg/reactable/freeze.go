package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) TryFreeze(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		//if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[ModifierCryo] < ZeroDur {
			return
		}
		consumed = r.triggerFreeze(r.Durability[ModifierCryo], a.Info.Durability)
		r.Durability[ModifierCryo] -= consumed
		r.Durability[ModifierCryo] = max(r.Durability[ModifierCryo], 0)
	case attributes.Cryo:
		if r.Durability[ModifierHydro] < ZeroDur {
			return
		}
		consumed := r.triggerFreeze(r.Durability[ModifierHydro], a.Info.Durability)
		r.Durability[ModifierHydro] -= consumed
		r.Durability[ModifierHydro] = max(r.Durability[ModifierHydro], 0)
	default:
		//should be here
		return
	}
	a.Reacted = true
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	r.core.Events.Emit(event.OnFrozen, r.self, a)
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
	if a.Info.StrikeType != combat.StrikeTypeBlunt || r.Durability[ModifierFrozen] < ZeroDur {
		return
	}
	//remove 200 freeze gauge if availabe
	r.Durability[ModifierFrozen] -= 200
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
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	ai.FlatDmg = 1.5 * calcReactionDmg(char, ai, em)
	//shatter is a self attack
	r.core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(r.self.Key()),
		-1,
		1,
	)

}

// add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b combat.Durability) combat.Durability {
	d := min(a, b)
	//trigger freeze should only addDurability and should not touch decay rate
	r.addDurability(ModifierFrozen, 2*d)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[ModifierFrozen] <= ZeroDur {
		r.Durability[ModifierFrozen] = 0
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
		r.core.QueueAttack(ai, combat.NewDefSingleTarget(r.self.Key()), -1, 1)
	}
}
