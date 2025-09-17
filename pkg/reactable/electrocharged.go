package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryAddEC(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// if there's still frozen left don't try to ec
	// game actively rejects ec reaction if frozen is present
	if r.Durability[info.ReactionModKeyFrozen] > info.ZeroDur {
		return false
	}

	// adding ec or hydro just adds to durability
	switch a.Info.Element {
	case attributes.Hydro:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur {
			return false
		}
		// add to hydro durability (can't add if the atk already reacted)
		// TODO: this shouldn't happen here
		if !a.Reacted {
			r.attachOrRefillNormalEle(info.ReactionModKeyHydro, a.Info.Durability)
		}
	case attributes.Electro:
		// if there's no existing hydro or electro then do nothing
		if r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
			return false
		}
		// add to electro durability (can't add if the atk already reacted)
		if !a.Reacted {
			r.attachOrRefillNormalEle(info.ReactionModKeyElectro, a.Info.Durability)
		}
	default:
		return false
	}

	a.Reacted = true
	r.core.Events.Emit(event.OnElectroCharged, r.self, a)

	// at this point ec is refereshed so we need to trigger a reaction
	// and change ownership
	atk := info.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeElectroCharged),
		AttackTag:        attacks.AttackTagECDamage,
		ICDTag:           attacks.ICDTagECDamage,
		ICDGroup:         attacks.ICDGroupReactionB,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, atk, em)
	atk.FlatDmg = 2.0 * flatdmg
	r.ecAtk = atk
	r.ecSnapshot = snap

	// if this is a new ec then trigger tick immediately and queue up ticks
	// otherwise do nothing
	// TODO: need to check if refresh ec triggers new tick immediately or not
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
		r.core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
			// target should be first, then snapshot
			n := args[0].(info.Target)
			a := args[1].(*info.AttackEvent)
			dmg := args[2].(float64)
			// TODO: there's no target index
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
			if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur || r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
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
	r.Durability[info.ReactionModKeyElectro] -= 10
	r.Durability[info.ReactionModKeyElectro] = max(0, r.Durability[info.ReactionModKeyElectro])
	r.Durability[info.ReactionModKeyHydro] -= 10
	r.Durability[info.ReactionModKeyHydro] = max(0, r.Durability[info.ReactionModKeyHydro])
	r.core.Log.NewEvent("ec wane",
		glog.LogElementEvent,
		-1,
	).
		Write("aura", "ec").
		Write("target", r.self.Key()).
		Write("hydro", r.Durability[info.ReactionModKeyHydro]).
		Write("electro", r.Durability[info.ReactionModKeyElectro])

	// ec is gone
	r.checkEC()
}

func (r *Reactable) checkEC() {
	if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur || r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
		r.ecTickSrc = -1
		r.core.Events.Unsubscribe(event.OnEnemyDamage, fmt.Sprintf("ec-%v", r.self.Key()))
		r.core.Log.NewEvent("ec expired",
			glog.LogElementEvent,
			-1,
		).
			Write("aura", "ec").
			Write("target", r.self.Key()).
			Write("hydro", r.Durability[info.ReactionModKeyHydro]).
			Write("electro", r.Durability[info.ReactionModKeyElectro])
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
		if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur || r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
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
