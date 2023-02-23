package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) TryOverload(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Electro:
		//must have pyro; pyro cant coexist (for now) so ok to ignore count?
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
			return false
		}
		//reduce; either gone or left; don't care how much actually reacted
		consumed = r.reduce(attributes.Pyro, a.Info.Durability, 1)
		r.burningCheck()
	case attributes.Pyro:
		//must have electro; gotta be careful with ec?
		if r.Durability[ModifierElectro] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	default:
		//should be here
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	//trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events.Emit(event.OnOverload, r.self, a)

	//trigger an overload attack
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(combat.Overload),
		AttackTag:        attacks.AttackTagOverloadDamage,
		ICDTag:           attacks.ICDTagOverloadDamage,
		ICDGroup:         combat.ICDGroupReactionB,
		StrikeType:       attacks.StrikeTypeBlunt,
		Element:          attributes.Pyro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := calcReactionDmg(char, atk, em)
	atk.FlatDmg = 2 * flatdmg
	r.core.QueueAttackWithSnap(atk, snap, combat.NewCircleHitOnTarget(r.self, nil, 3), 1)

	return true
}
