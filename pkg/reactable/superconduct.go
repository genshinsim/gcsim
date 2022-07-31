package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) trySuperconduct(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for non frozen one
	if r.Durability[attributes.Frozen] >= ZeroDur {
		return
	}
	switch a.Info.Element {
	case attributes.Electro:
		if r.Durability[attributes.Cryo] < ZeroDur {
			return
		}
		r.reduce(attributes.Cryo, a.Info.Durability, 1)
		a.Info.Durability = 0
	case attributes.Cryo:
		//could be ec potentially
		if r.Durability[attributes.Electro] < ZeroDur {
			return
		}
		rd := r.reduce(attributes.Electro, a.Info.Durability, 1)
		//if there's hydro as well then don't consume all the durability
		if r.Durability[attributes.Hydro] > ZeroDur {
			a.Info.Durability -= rd
		} else {
			a.Info.Durability = 0
		}
	default:
		return
	}

	r.queueSuperconduct(a)

}

func (r *Reactable) tryFrozenSuperconduct(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for frozen
	if r.Durability[attributes.Frozen] < ZeroDur {
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
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Superconduct),
		AttackTag:        combat.AttackTagSuperconductDamage,
		ICDTag:           combat.ICDTagSuperconductDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	atk.FlatDmg = 0.5 * r.calcReactionDmg(atk, em)
	r.core.QueueAttack(atk, combat.NewCircleHit(r.self, 3, true, combat.TargettableEnemy), -1, 1)
}
