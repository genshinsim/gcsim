package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (r *Reactable) tryAddEC(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	//adding ec or hydro just adds to durability
	switch a.Info.Element {
	case attributes.Hydro:
		//if there's no existing hydro or electro then do nothing
		if r.Durability[attributes.Electro] < ZeroDur {
			return
		}
		if r.Durability[attributes.Hydro] < ZeroDur {
			//attach
			r.tryAttach(attributes.Hydro, &a.Info.Durability)
		} else {
			r.tryRefill(attributes.Hydro, &a.Info.Durability)
		}
		//add to hydro durability
	case attributes.Electro:
		//if there's no existing hydro or electro then do nothing
		if r.Durability[attributes.Hydro] < ZeroDur {
			return
		}
		//add to electro durability
		if r.Durability[attributes.Electro] < ZeroDur {
			//attach
			r.tryAttach(attributes.Electro, &a.Info.Durability)
		} else {
			r.tryRefill(attributes.Electro, &a.Info.Durability)
		}
	default:
		return
	}

	r.core.Events.Emit(event.OnElectroCharged, r.self, a)

	//at this point ec is refereshed so we need to trigger a reaction
	//and change ownership
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.ElectroCharged),
		AttackTag:        combat.AttackTagECDamage,
		ICDTag:           combat.ICDTagECDamage,
		ICDGroup:         combat.ICDGroupReactionB,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
	}
	em := r.core.CharAttr.Stat(a.Info.ActorIndex, attributes.EM)
	atk.FlatDmg = 1.2 * r.calcReactionDmg(atk, em)
	r.ecSnapshot = atk

	//if this is a new ec then trigger tick immediately and queue up ticks
	//otherwise do nothing
	//TODO: need to check if refresh ec triggers new tick immediately or not
	if r.ecTickSrc == -1 {
		r.ecTickSrc = r.core.F

		r.core.QueueAttack(
			r.ecSnapshot,
			combat.NewDefSingleTarget(r.self.Index(), r.self.Type()),
			-1,
			10,
		)

		r.core.Tasks.Add(r.nextTick(r.core.F), 60+10)
		//subscribe to wane ticks
		r.core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
			//target should be first, then snapshot
			n := args[0].(combat.Target)
			a := args[1].(*combat.AttackEvent)
			dmg := args[2].(float64)
			//TODO: there's no target index
			if n.Index() != r.self.Index() {
				return false
			}
			if a.Info.AttackTag != combat.AttackTagECDamage {
				return false
			}
			//ignore if this dmg instance has been wiped out due to icd
			if dmg == 0 {
				return false
			}
			//ignore if we no longer have both electro and hydro
			if r.Durability[attributes.Electro] < ZeroDur || r.Durability[attributes.Hydro] < ZeroDur {
				return true
			}

			//wane in 0.1 seconds
			r.core.Tasks.Add(func() {
				r.waneEC()
			}, 6)
			return false
		}, fmt.Sprintf("ec-%v", r.self.Index()))
	}

	//ticks are 60 frames since last tick
	//taking tick dmg resets last tick
}

func (r *Reactable) waneEC() {
	r.Durability[attributes.Electro] -= 10
	r.Durability[attributes.Electro] = max(0, r.Durability[attributes.Electro])
	r.Durability[attributes.Hydro] -= 10
	r.Durability[attributes.Hydro] = max(0, r.Durability[attributes.Hydro])
	r.core.Log.NewEvent("ec wane",
		glog.LogElementEvent,
		-1,
		"aura", "ec",
		"target", r.self.Index(),
		"hydro", r.Durability[attributes.Hydro],
		"electro", r.Durability[attributes.Electro],
	)
	//ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[attributes.Electro] < ZeroDur || r.Durability[attributes.Hydro] < ZeroDur {
		r.ecTickSrc = -1
		r.core.Events.Unsubscribe(event.OnDamage, fmt.Sprintf("ec-%v", r.self.Index()))
		r.core.Log.NewEvent("ec expired",
			glog.LogElementEvent,
			-1,
			"aura", "ec",
			"target", r.self.Index(),
			"hydro", r.Durability[attributes.Hydro],
			"electro", r.Durability[attributes.Electro],
		)
	}
}

func (r *Reactable) nextTick(src int) func() {
	return func() {
		if r.ecTickSrc != src {
			//source changed, do nothing
			return
		}
		//ec SHOULD be active still, since if not we would have
		//called cleanup and set source to -1
		if r.Durability[attributes.Electro] < ZeroDur || r.Durability[attributes.Hydro] < ZeroDur {
			return
		}

		//so ec is active, which means both aura must still have value > 0; so we can do dmg
		r.core.QueueAttack(
			r.ecSnapshot,
			combat.NewDefSingleTarget(r.self.Index(), r.self.Type()),
			-1,
			0,
		)

		//queue up next tick
		r.core.Tasks.Add(r.nextTick(src), 60)
	}
}
