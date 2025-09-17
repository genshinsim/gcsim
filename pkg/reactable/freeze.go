package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryFreeze(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// so if already frozen there are 2 cases:
	// 1. src exists but no other coexisting -> attach
	// 2. src does not exist but opposite coexists -> add to freeze durability
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		// if cryo exists we'll trigger freeze regardless if frozen already coexists
		if r.Durability[info.ReactionModKeyCryo] < info.ZeroDur {
			return false
		}
		consumed = r.triggerFreeze(r.Durability[info.ReactionModKeyCryo], a.Info.Durability)
		r.Durability[info.ReactionModKeyCryo] -= consumed
		r.Durability[info.ReactionModKeyCryo] = max(r.Durability[info.ReactionModKeyCryo], 0)
	case attributes.Cryo:
		if r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
			return false
		}
		consumed := r.triggerFreeze(r.Durability[info.ReactionModKeyHydro], a.Info.Durability)
		r.Durability[info.ReactionModKeyHydro] -= consumed
		r.Durability[info.ReactionModKeyHydro] = max(r.Durability[info.ReactionModKeyHydro], 0)
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

func (r *Reactable) PoiseDMGCheck(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyFrozen] < info.ZeroDur {
		return false
	}
	if a.Info.StrikeType != attacks.StrikeTypeBlunt {
		return false
	}
	// remove frozen durability according to poise dmg
	r.Durability[info.ReactionModKeyFrozen] -= info.Durability(0.15 * a.Info.PoiseDMG)
	r.checkFreeze()
	return true
}

func (r *Reactable) ShatterCheck(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyFrozen] < info.ZeroDur {
		return false
	}
	if a.Info.StrikeType != attacks.StrikeTypeBlunt && a.Info.Element != attributes.Geo {
		return false
	}
	// remove 200 freeze gauge if available
	r.Durability[info.ReactionModKeyFrozen] -= 200
	r.checkFreeze()

	r.core.Events.Emit(event.OnShatter, r.self, a)

	// 0.2s gcd on shatter attack
	if r.shatterGCD == -1 || r.core.F >= r.shatterGCD {
		r.shatterGCD = r.core.F + 0.2*60
		// trigger shatter attack
		ai := info.AttackInfo{
			ActorIndex:       a.Info.ActorIndex,
			DamageSrc:        r.self.Key(),
			Abil:             string(info.ReactionTypeShatter),
			AttackTag:        attacks.AttackTagShatter,
			ICDTag:           attacks.ICDTagShatter,
			ICDGroup:         attacks.ICDGroupReactionA,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Physical,
			IgnoreDefPercent: 1,
		}
		char := r.core.Player.ByIndex(a.Info.ActorIndex)
		em := char.Stat(attributes.EM)
		flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, ai, em)
		ai.FlatDmg = 3.0 * flatdmg
		// shatter is a self attack
		r.core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewSingleTargetHit(r.self.Key()),
			0,
		)
	}

	return true
}

// add to freeze durability and return amount of durability consumed
func (r *Reactable) triggerFreeze(a, b info.Durability) info.Durability {
	d := min(a, b)
	if r.FreezeResist >= 1 {
		return d
	}
	// trigger freeze should only addDurability and should not touch decay rate
	r.attachOverlap(info.ReactionModKeyFrozen, 2*d, info.ZeroDur)
	return d
}

func (r *Reactable) checkFreeze() {
	if r.Durability[info.ReactionModKeyFrozen] <= info.ZeroDur {
		r.Durability[info.ReactionModKeyFrozen] = 0
		r.core.Events.Emit(event.OnAuraDurabilityDepleted, r.self, attributes.Frozen)
		// trigger another attack here, purely for the purpose of breaking bubbles >.>
		ai := info.AttackInfo{
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
		// TODO: delay attack by 1 frame ok?
		r.core.QueueAttack(ai, combat.NewSingleTargetHit(r.self.Key()), -1, 0)
	}
}
