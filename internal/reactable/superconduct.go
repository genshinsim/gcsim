package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (r *Reactable) trySuperconduct(a *coretype.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for non frozen one
	if r.Durability[coretype.Frozen] >= ZeroDur {
		return
	}
	switch a.Info.Element {
	case core.Electro:
		if r.Durability[coretype.Cryo] < ZeroDur {
			return
		}
		r.reduce(coretype.Cryo, a.Info.Durability, 1)
		a.Info.Durability = 0
	case coretype.Cryo:
		//could be ec potentially
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
		return
	}

	r.queueSuperconduct(a)

}

func (r *Reactable) tryFrozenSuperconduct(a *coretype.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//this is for frozen
	if r.Durability[coretype.Frozen] < ZeroDur {
		return
	}
	//
	switch a.Info.Element {
	case core.Electro:
		//TODO: the assumption here is we first reduce cryo, and if there's any
		//src durability left, we reduce frozen. note that it's still only one
		//superconduct reaction
		a.Info.Durability -= r.reduce(coretype.Cryo, a.Info.Durability, 1)
		r.reduce(coretype.Frozen, a.Info.Durability, 1)
		a.Info.Durability = 0
	default:
		return
	}

	r.queueSuperconduct(a)

}

func (r *Reactable) queueSuperconduct(a *coretype.AttackEvent) {
	r.core.Events.Emit(core.OnSuperconduct, r.self, a)

	//superconduct attack
	atk := core.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(core.Superconduct),
		AttackTag:        core.AttackTagSuperconductDamage,
		ICDTag:           core.ICDTagSuperconductDamage,
		ICDGroup:         core.ICDGroupReactionA,
		Element:          coretype.Cryo,
		IgnoreDefPercent: 1,
	}
	atk.FlatDmg = 0.5 * r.calcReactionDmg(atk)
	r.core.Combat.QueueAttack(atk, core.NewDefCircHit(3, true, coretype.TargettableEnemy), -1, 1, superconductPhysShred)
}

func superconductPhysShred(a core.AttackCB) {
	a.Target.AddResMod("superconductphysshred", core.ResistMod{
		Duration: 12 * 60,
		Ele:      core.Physical,
		Value:    -0.4,
	})
}
