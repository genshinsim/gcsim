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
		if r.Durability[ModifierElectro] < ZeroDur {
			return
		}
		//add to hydro durability
		//TODO: this shouldn't happen here
		r.attachOrRefillNormalEle(ModifierHydro, a.Info.Durability)
	case attributes.Electro:
		//if there's no existing hydro or electro then do nothing
		if r.Durability[ModifierHydro] < ZeroDur {
			return
		}
		//add to electro durability
		r.attachOrRefillNormalEle(ModifierElectro, a.Info.Durability)
	default:
		return
	}

	a.Reacted = true
	r.core.Events.Emit(event.OnElectroCharged, r.self, a)

	//at this point ec is refereshed so we need to trigger a reaction
	//and change ownership
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(combat.ElectroCharged),
		AttackTag:        combat.AttackTagECDamage,
		ICDTag:           combat.ICDTagECDamage,
		ICDGroup:         combat.ICDGroupReactionB,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	atk.FlatDmg = 1.2 * calcReactionDmg(char, atk, em)
	r.ecSnapshot = atk

	//if this is a new ec then trigger tick immediately and queue up ticks
	//otherwise do nothing
	//TODO: need to check if refresh ec triggers new tick immediately or not
	if r.ecTickSrc == -1 {
		r.ecTickSrc = r.core.F

		r.core.QueueAttack(
			r.ecSnapshot,
			combat.NewDefSingleTarget(r.self.Key(), r.self.Type()),
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
			if r.Durability[ModifierElectro] < ZeroDur || r.Durability[ModifierHydro] < ZeroDur {
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
	r.Durability[ModifierElectro] -= 10
	r.Durability[ModifierElectro] = max(0, r.Durability[ModifierElectro])
	r.Durability[ModifierHydro] -= 10
	r.Durability[ModifierHydro] = max(0, r.Durability[ModifierHydro])
	r.core.Log.NewEvent("ec wane",
		glog.LogElementEvent,
		-1,
	).
		Write("aura", "ec").
		Write("target", r.self.Index()).
		Write("hydro", r.Durability[ModifierHydro]).
		Write("electro", r.Durability[ModifierElectro])

	//ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[ModifierElectro] < ZeroDur || r.Durability[ModifierHydro] < ZeroDur {
		r.ecTickSrc = -1
		r.core.Events.Unsubscribe(event.OnDamage, fmt.Sprintf("ec-%v", r.self.Index()))
		r.core.Log.NewEvent("ec expired",
			glog.LogElementEvent,
			-1,
		).
			Write("aura", "ec").
			Write("target", r.self.Index()).
			Write("hydro", r.Durability[ModifierHydro]).
			Write("electro", r.Durability[ModifierElectro])

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
		if r.Durability[ModifierElectro] < ZeroDur || r.Durability[ModifierHydro] < ZeroDur {
			return
		}

		//so ec is active, which means both aura must still have value > 0; so we can do dmg
		r.core.QueueAttack(
			r.ecSnapshot,
			combat.NewDefSingleTarget(r.self.Key(), r.self.Type()),
			-1,
			0,
		)

		//queue up next tick
		r.core.Tasks.Add(r.nextTick(src), 60)
	}
}
