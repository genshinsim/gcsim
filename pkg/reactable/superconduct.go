package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) TrySuperconduct(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for non frozen one
	if r.Durability[ModifierFrozen] >= ZeroDur {
		return
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Electro:
		if r.Durability[ModifierCryo] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 1)
	case attributes.Cryo:
		//could be ec potentially
		if r.Durability[ModifierElectro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	default:
		return
	}

	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.queueSuperconduct(a)
}

func (r *Reactable) TryFrozenSuperconduct(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for frozen
	if r.Durability[ModifierFrozen] < ZeroDur {
		return
	}
	switch a.Info.Element {
	case attributes.Electro:
		//TODO: the assumption here is we first reduce cryo, and if there's any
		//src durability left, we reduce frozen. note that it's still only one
		//superconduct reaction
		a.Info.Durability -= r.reduce(attributes.Cryo, a.Info.Durability, 1)
		r.reduce(attributes.Frozen, a.Info.Durability, 1)
		a.Info.Durability = 0
		a.Reacted = true
	default:
		return
	}

	r.queueSuperconduct(a)

}

func (r *Reactable) queueSuperconduct(a *combat.AttackEvent) {
	r.core.Events.Emit(event.OnSuperconduct, r.self, a)

	//superconduct attack
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(combat.Superconduct),
		AttackTag:        combat.AttackTagSuperconductDamage,
		ICDTag:           combat.ICDTagSuperconductDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	atk.FlatDmg = 0.5 * calcReactionDmg(char, atk, em)
	r.core.QueueAttack(atk, combat.NewCircleHit(r.self, 3), -1, 1)
}
