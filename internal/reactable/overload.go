package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (r *Reactable) tryOverload(a *core.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	switch a.Info.Element {
	case core.Electro:
		//must have pyro; pyro cant coexist (for now) so ok to ignore count?
		if r.Durability[core.Pyro] < ZeroDur {
			return
		}
		//reduce; either gone or left; don't care how much actually reacted
		r.reduce(core.Pyro, a.Info.Durability, 1)
		//since there's nothing else to react with, reduce durability to 0
		a.Info.Durability = 0
	case core.Pyro:
		//must have electro; gotta be careful with ec?
		if r.Durability[core.Electro] < ZeroDur {
			return
		}
		rd := r.reduce(core.Electro, a.Info.Durability, 1)
		//if there's hydro as well then don't consume all the durability
		if r.Durability[core.Hydro] > ZeroDur {
			a.Info.Durability -= rd
		} else {
			a.Info.Durability = 0
		}
	default:
		//should be here
		return
	}

	//trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events.Emit(core.OnOverload, r.self, a)

	//trigger an overload attack
	atk := core.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(core.Overload),
		AttackTag:        core.AttackTagOverloadDamage,
		ICDTag:           core.ICDTagOverloadDamage,
		ICDGroup:         core.ICDGroupReactionB,
		StrikeType:       core.StrikeTypeBlunt,
		Element:          core.Pyro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Chars[a.Info.ActorIndex]
	em := char.Stat(core.EM)
	atk.FlatDmg = 2 * r.calcReactionDmg(atk, em)
	r.core.Combat.QueueAttack(atk, core.NewDefCircHit(3, true, core.TargettableEnemy), -1, 1)
}
