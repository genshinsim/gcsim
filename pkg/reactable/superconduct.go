package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TrySuperconduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for non frozen one
	if r.Durability[info.ReactionModKeyFrozen] >= info.ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Electro:
		if r.Durability[info.ReactionModKeyCryo] < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 1)
	case attributes.Cryo:
		// could be ec potentially
		if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	default:
		return false
	}

	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.queueSuperconduct(a)
	return true
}

func (r *Reactable) TryFrozenSuperconduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for frozen
	if r.Durability[info.ReactionModKeyFrozen] < info.ZeroDur {
		return false
	}
	switch a.Info.Element {
	case attributes.Electro:
		// TODO: the assumption here is we first reduce cryo, and if there's any
		// src durability left, we reduce frozen. note that it's still only one
		// superconduct reaction
		a.Info.Durability -= r.reduce(attributes.Cryo, a.Info.Durability, 1)
		r.reduce(attributes.Frozen, a.Info.Durability, 1)
		a.Info.Durability = 0
		a.Reacted = true
	default:
		return false
	}

	r.queueSuperconduct(a)

	return false
}

func (r *Reactable) queueSuperconduct(a *info.AttackEvent) {
	r.core.Events.Emit(event.OnSuperconduct, r.self, a)

	// 0.1s gcd on superconduct attack
	if r.superconductGCD != -1 && r.core.F < r.superconductGCD {
		return
	}
	r.superconductGCD = r.core.F + 0.1*60

	// superconduct attack
	atk := info.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeSuperconduct),
		AttackTag:        attacks.AttackTagSuperconductDamage,
		ICDTag:           attacks.ICDTagSuperconductDamage,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, atk, em)
	atk.FlatDmg = 1.5 * flatdmg
	r.core.QueueAttackWithSnap(atk, snap, combat.NewCircleHitOnTarget(r.self, nil, 3), 1)
}
