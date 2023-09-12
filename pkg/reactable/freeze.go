package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

func (r *Reactable) TryFreeze(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	var consumed reactions.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		// if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[Cryo] < ZeroDur {
			return false
		}
		consumed = r.triggerFreeze(r.Durability[Cryo], a.Info.Durability)
		r.Durability[Cryo] -= consumed
		r.Durability[Cryo] = max(r.Durability[Cryo], 0)
	case attributes.Cryo:
		if r.Durability[Hydro] < ZeroDur {
			return false
		}
		consumed := r.triggerFreeze(r.Durability[Hydro], a.Info.Durability)
		r.Durability[Hydro] -= consumed
		r.Durability[Hydro] = max(r.Durability[Hydro], 0)
	default:
		// should be here
		return false
	}
	a.Reacted = true
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	r.core.Events.Emit(event.OnFrozen, r.self, a)
	return true
}

func max(a, b reactions.Durability) reactions.Durability {
	if a > b {
		return a
	}
	return b
}

func min(a, b reactions.Durability) reactions.Durability {
	if a > b {
		return b
	}
	return a
}

func (r *Reactable) ShatterCheck(a *combat.AttackEvent) bool {
	if r.Durability[Frozen] < ZeroDur {
		return false
	}
	if a.Info.StrikeType != attacks.StrikeTypeBlunt && a.Info.Element != attributes.Geo {
		return false
	}
	// remove 200 freeze gauge if availabe
	r.Durability[Frozen] -= 200
	r.checkFreeze()
	// trigger shatter attack
	r.core.Events.Emit(event.OnShatter, r.self, a)
	ai := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(reactions.Shatter),
		AttackTag:        attacks.AttackTagShatter,
		ICDTag:           attacks.ICDTagShatter,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Physical,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := calcReactionDmg(char, ai, em)
	ai.FlatDmg = 1.5 * flatdmg
	// shatter is a self attack
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewSingleTargetHit(r.self.Key()),
		1,
	)
	return true
}

// add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b reactions.Durability) reactions.Durability {
	d := min(a, b)
	if r.FreezeResist >= 1 {
		return d
	}
	// trigger freeze should only addDurability and should not touch decay rate
	r.attachOverlap(Frozen, 2*d, ZeroDur)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[Frozen] <= ZeroDur {
		r.Durability[Frozen] = 0
		r.core.Events.Emit(event.OnAuraDurabilityDepleted, r.self, attributes.Frozen)
		// trigger another attack here, purely for the purpose of breaking bubbles >.>
		ai := combat.AttackInfo{
			ActorIndex:  0,
			DamageSrc:   r.self.Key(),
			Abil:        "Freeze Broken",
			AttackTag:   attacks.AttackTagNone,
			ICDTag:      attacks.ICDTagNone,
			ICDGroup:    attacks.ICDGroupDefault,
			StrikeType:  attacks.StrikeTypeDefault,
			Element:     attributes.NoElement,
			SourceIsSim: true,
			DoNotLog:    true,
		}
		//TODO: delay attack by 1 frame ok?
		r.core.QueueAttack(ai, combat.NewSingleTargetHit(r.self.Key()), -1, 1)
	}
}
