package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/template/crystallize"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryCrystallizeElectro(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyElectro] > info.ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Electro, info.ReactionTypeCrystallizeElectro, event.OnCrystallizeElectro)
	}
	return false
}

func (r *Reactable) TryCrystallizeHydro(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyHydro] > info.ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Hydro, info.ReactionTypeCrystallizeHydro, event.OnCrystallizeHydro)
	}
	return false
}

func (r *Reactable) TryCrystallizeCryo(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyCryo] > info.ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Cryo, info.ReactionTypeCrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

func (r *Reactable) TryCrystallizePyro(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyPyro] > info.ZeroDur || r.Durability[info.ReactionModKeyBurning] > info.ZeroDur {
		reacted := r.tryCrystallizeWithEle(a, attributes.Pyro, info.ReactionTypeCrystallizePyro, event.OnCrystallizePyro)
		r.burningCheck()
		return reacted
	}
	return false
}

func (r *Reactable) TryCrystallizeFrozen(a *info.AttackEvent) bool {
	if r.Durability[info.ReactionModKeyFrozen] > info.ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Frozen, info.ReactionTypeCrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

func (r *Reactable) tryCrystallizeWithEle(a *info.AttackEvent, ele attributes.Element, rt info.ReactionType, evt event.Event) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	if r.crystallizeGCD != -1 && r.core.F < r.crystallizeGCD {
		return false
	}
	r.crystallizeGCD = r.core.F + 60
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	r.addCrystallizeShard(char, rt, ele, r.core.F)
	// reduce
	r.reduce(ele, a.Info.Durability, 0.5)
	// TODO: confirm u can only crystallize once
	a.Info.Durability = 0
	a.Reacted = true
	// event
	r.core.Events.Emit(evt, r.self, a)
	// check freeze + ec
	switch {
	case ele == attributes.Electro && r.Durability[info.ReactionModKeyHydro] > info.ZeroDur:
		r.checkEC()
	case ele == attributes.Frozen:
		r.checkFreeze()
	}
	return true
}

type crystalizeChar interface {
	Snapshot(*info.AttackInfo) info.Snapshot
	Index() int
}

func (r *Reactable) addCrystallizeShard(char crystalizeChar, rt info.ReactionType, typ attributes.Element, src int) {
	// delay shard spawn
	r.core.Tasks.Add(func() {
		// grab current snapshot for shield
		ai := info.AttackInfo{
			ActorIndex: char.Index(),
			DamageSrc:  r.self.Key(),
			Abil:       string(rt),
		}
		snap := char.Snapshot(&ai)
		lvl := snap.CharLvl
		// shield snapshots em on shard spawn
		em := snap.Stats[attributes.EM]
		// expiry will get set properly later
		shd := crystallize.NewShield(char.Index(), typ, src, lvl, em, -1)
		cs := crystallize.NewShard(r.core, r.self.Shape(), shd)
		r.core.Combat.AddGadget(cs)
		r.core.Log.NewEvent(
			fmt.Sprintf("%v crystallize shard spawned", cs.Shield.Ele),
			glog.LogElementEvent,
			cs.Shield.ActorIndex,
		).
			Write("src", cs.Src()).
			Write("expiry", cs.Expiry()).
			Write("earliest_pickup", cs.EarliestPickup)
	}, 23)
}
