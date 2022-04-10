package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (r *Reactable) tryAddEC(a *core.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	//adding ec or hydro just adds to durability
	switch a.Info.Element {
	case core.Hydro:
		//if there's no existing hydro or electro then do nothing
		if r.Durability[core.Electro] < ZeroDur {
			return
		}
		if r.Durability[core.Hydro] < ZeroDur {
			//attach
			r.tryAttach(core.Hydro, &a.Info.Durability)
		} else {
			r.tryRefill(core.Hydro, &a.Info.Durability)
		}
		//add to hydro durability
	case core.Electro:
		//if there's no existing hydro or electro then do nothing
		if r.Durability[core.Hydro] < ZeroDur {
			return
		}
		//add to electro durability
		if r.Durability[core.Electro] < ZeroDur {
			//attach
			r.tryAttach(core.Electro, &a.Info.Durability)
		} else {
			r.tryRefill(core.Electro, &a.Info.Durability)
		}
	default:
		return
	}

	r.core.Events.Emit(core.OnElectroCharged, r.self, a)

	//at this point ec is refereshed so we need to trigger a reaction
	//and change ownership
	atk := core.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(core.ElectroCharged),
		AttackTag:        core.AttackTagECDamage,
		ICDTag:           core.ICDTagECDamage,
		ICDGroup:         core.ICDGroupReactionB,
		Element:          core.Electro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Chars[a.Info.ActorIndex]
	em := char.Stat(core.EM)
	atk.FlatDmg = 1.2 * r.calcReactionDmg(atk, em)
	r.ecSnapshot = atk

	//if this is a new ec then trigger tick immediately and queue up ticks
	//otherwise do nothing
	//TODO: need to check if refresh ec triggers new tick immediately or not
	if r.ecTickSrc == -1 {
		r.ecTickSrc = r.core.F

		r.core.Combat.QueueAttack(
			r.ecSnapshot,
			core.NewDefSingleTarget(r.self.Index(), r.self.Type()),
			-1,
			10,
		)

		r.core.Tasks.Add(r.nextTick(r.core.F), 60+10)
		//subscribe to wane ticks
		r.core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			//target should be first, then snapshot
			n := args[0].(core.Target)
			a := args[1].(*core.AttackEvent)
			dmg := args[2].(float64)
			//TODO: there's no target index
			if n.Index() != r.self.Index() {
				return false
			}
			if a.Info.AttackTag != core.AttackTagECDamage {
				return false
			}
			//ignore if this dmg instance has been wiped out due to icd
			if dmg == 0 {
				return false
			}
			//ignore if we no longer have both electro and hydro
			if r.Durability[core.Electro] < ZeroDur || r.Durability[core.Hydro] < ZeroDur {
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
	r.Durability[core.Electro] -= 10
	r.Durability[core.Electro] = max(0, r.Durability[core.Electro])
	r.Durability[core.Hydro] -= 10
	r.Durability[core.Hydro] = max(0, r.Durability[core.Hydro])
	r.core.Log.NewEvent("ec wane",
		core.LogElementEvent,
		-1,
		"aura", "ec",
		"target", r.self.Index(),
		"hydro", r.Durability[core.Hydro],
		"electro", r.Durability[core.Electro],
	)
	//ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[core.Electro] < ZeroDur || r.Durability[core.Hydro] < ZeroDur {
		r.ecTickSrc = -1
		r.core.Events.Unsubscribe(core.OnDamage, fmt.Sprintf("ec-%v", r.self.Index()))
		r.core.Log.NewEvent("ec expired",
			core.LogElementEvent,
			-1,
			"aura", "ec",
			"target", r.self.Index(),
			"hydro", r.Durability[core.Hydro],
			"electro", r.Durability[core.Electro],
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
		if r.Durability[core.Electro] < ZeroDur || r.Durability[core.Hydro] < ZeroDur {
			return
		}

		//so ec is active, which means both aura must still have value > 0; so we can do dmg
		r.core.Combat.QueueAttack(
			r.ecSnapshot,
			core.NewDefSingleTarget(r.self.Index(), r.self.Type()),
			-1,
			0,
		)

		//queue up next tick
		r.core.Tasks.Add(r.nextTick(src), 60)
	}
}
