package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryOverload(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	switch a.Info.Element {
	case attributes.Electro:
		//must have pyro; pyro cant coexist (for now) so ok to ignore count?
		if r.Durability[attributes.Pyro] < ZeroDur {
			return
		}
		//reduce; either gone or left; don't care how much actually reacted
		r.reduce(attributes.Pyro, a.Info.Durability, 1)
		//since there's nothing else to react with, reduce durability to 0
		a.Info.Durability = 0
	case attributes.Pyro:
		//must have electro; gotta be careful with ec?
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
		//should be here
		return
	}

	//trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events.Emit(event.OnOverload, r.self, a)

	//trigger an overload attack
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Overload),
		AttackTag:        combat.AttackTagOverloadDamage,
		ICDTag:           combat.ICDTagOverloadDamage,
		ICDGroup:         combat.ICDGroupReactionB,
		StrikeType:       combat.StrikeTypeBlunt,
		Element:          attributes.Pyro,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.Stat(a.Info.ActorIndex, attributes.EM)
	atk.FlatDmg = 2 * r.calcReactionDmg(atk, em)
	r.core.QueueAttack(atk, combat.NewDefCircHit(3, true, combat.TargettableEnemy), -1, 1)
}
