package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

func (r *Reactable) TryAddEC(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// if there's still frozen left don't try to ec
	// game actively rejects ec reaction if frozen is present
	if r.Durability[Frozen] > ZeroDur {
		return false
	}

	// adding ec or hydro just adds to durability
	switch a.Info.Element {
	case attributes.Hydro:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[Electro] < ZeroDur {
			return false
		}
		// add to hydro durability (can't add if the atk already reacted)
		//TODO: this shouldn't happen here
		if !a.Reacted {
			r.attachOrRefillNormalEle(Hydro, a.Info.Durability)
		}
	case attributes.Electro:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[Hydro] < ZeroDur {
			return false
		}
		// add to electro durability (can't add if the atk already reacted)
		if !a.Reacted {
			r.attachOrRefillNormalEle(Electro, a.Info.Durability)
		}
	default:
		return false
	}

	a.Reacted = true
	r.core.Events.Emit(event.OnElectroCharged, r.self, a)

	// at this point ec is refereshed so we need to trigger a reaction
	// and change ownership
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(reactions.ElectroCharged),
		AttackTag:        attacks.AttackTagECDamage,
		ICDTag:           attacks.ICDTagECDamage,
		ICDGroup:         attacks.ICDGroupReactionB,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := calcReactionDmg(char, atk, em)
	atk.FlatDmg = 2.0 * flatdmg
	r.ecAtk = atk
	r.ecSnapshot = snap

	// if this is a new ec then trigger tick immediately and queue up ticks
	// otherwise do nothing
	//TODO: need to check if refresh ec triggers new tick immediately or not
	if r.ecTickSrc == -1 {
		r.ecTickSrc = r.core.F
		r.core.QueueAttackWithSnap(
			r.ecAtk,
			r.ecSnapshot,
			combat.NewSingleTargetHit(r.self.Key()),
			10,
		)

		r.core.Tasks.Add(r.nextTick(r.core.F), 60+10)
		// subscribe to wane ticks
		r.core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			// target should be first, then snapshot
			n := args[0].(combat.Target)
			a := args[1].(*combat.AttackEvent)
			dmg := args[2].(float64)
			//TODO: there's no target index
			if n.Key() != r.self.Key() {
				return false
			}
			if a.Info.AttackTag != attacks.AttackTagECDamage {
				return false
			}
			// ignore if this dmg instance has been wiped out due to icd
			if dmg == 0 {
				return false
			}
			// ignore if we no longer have both electro and hydro
			if r.Durability[Electro] < ZeroDur || r.Durability[Hydro] < ZeroDur {
				return true
			}

			// wane in 0.1 seconds
			r.core.Tasks.Add(func() {
				r.waneEC()
			}, 6)
			return false
		}, fmt.Sprintf("ec-%v", r.self.Key()))
	}

	// ticks are 60 frames since last tick
	// taking tick dmg resets last tick
	return true
}

func (r *Reactable) waneEC() {
	r.Durability[Electro] -= 10
	r.Durability[Electro] = max(0, r.Durability[Electro])
	r.Durability[Hydro] -= 10
	r.Durability[Hydro] = max(0, r.Durability[Hydro])
	r.core.Log.NewEvent("ec wane",
		glog.LogElementEvent,
		-1,
	).
		Write("aura", "ec").
		Write("target", r.self.Key()).
		Write("hydro", r.Durability[Hydro]).
		Write("electro", r.Durability[Electro])

	// ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[Electro] < ZeroDur || r.Durability[Hydro] < ZeroDur {
		r.ecTickSrc = -1
		r.core.Events.Unsubscribe(event.OnEnemyDamage, fmt.Sprintf("ec-%v", r.self.Key()))
		r.core.Log.NewEvent("ec expired",
			glog.LogElementEvent,
			-1,
		).
			Write("aura", "ec").
			Write("target", r.self.Key()).
			Write("hydro", r.Durability[Hydro]).
			Write("electro", r.Durability[Electro])
	}
}

func (r *Reactable) nextTick(src int) func() {
	return func() {
		if r.ecTickSrc != src {
			// source changed, do nothing
			return
		}
		// ec SHOULD be active still, since if not we would have
		// called cleanup and set source to -1
		if r.Durability[Electro] < ZeroDur || r.Durability[Hydro] < ZeroDur {
			return
		}

		// so ec is active, which means both aura must still have value > 0; so we can do dmg
		r.core.QueueAttackWithSnap(
			r.ecAtk,
			r.ecSnapshot,
			combat.NewSingleTargetHit(r.self.Key()),
			0,
		)

		// queue up next tick
		r.core.Tasks.Add(r.nextTick(src), 60)
	}
}
