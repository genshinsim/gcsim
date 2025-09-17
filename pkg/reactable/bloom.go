package reactable

import (
	"github.com/genshinsim/gcsim/internal/template/dendrocore"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryBloom(a *info.AttackEvent) bool {
	// can be hydro bloom, dendro bloom, or quicken bloom
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		// this part is annoying. bloom will happen if any of the dendro like aura is present
		// so we gotta check for all 3...
		switch {
		case r.Durability[info.ReactionModKeyDendro] > info.ZeroDur:
		case r.Durability[info.ReactionModKeyQuicken] > info.ZeroDur:
		case r.Durability[info.ReactionModKeyBurningFuel] > info.ZeroDur:
		default:
			return false
		}
		// reduce only check for one element so have to call twice to check for quicken as well
		consumed = r.reduce(attributes.Dendro, a.Info.Durability, 0.5)
		f := r.reduce(attributes.Quicken, a.Info.Durability, 0.5)
		if f > consumed {
			consumed = f
		}
	case attributes.Dendro:
		if r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Hydro, a.Info.Durability, 2)
	default:
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.addBloomGadget(a)
	r.core.Events.Emit(event.OnBloom, r.self, a)
	return true
}

// this function should only be called after a catalyze reaction (queued to the end of current frame)
// this reaction will check if any hydro exists and if so trigger a bloom reaction
func (r *Reactable) tryQuickenBloom(a *info.AttackEvent) {
	if r.Durability[info.ReactionModKeyQuicken] < info.ZeroDur {
		// this should be a sanity check; should not happen realistically unless something wipes off
		// the quicken immediately (same frame) after catalyze
		return
	}
	if r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
		return
	}
	avail := r.Durability[info.ReactionModKeyQuicken]
	consumed := r.reduce(attributes.Hydro, avail, 2)
	r.Durability[info.ReactionModKeyQuicken] -= consumed

	r.addBloomGadget(a)
	r.core.Events.Emit(event.OnBloom, r.self, a)
}

func (r *Reactable) addBloomGadget(a *info.AttackEvent) {
	r.core.Tasks.Add(func() {
		t := dendrocore.New(r.core, r.self.Shape(), a)
		r.core.Combat.AddGadget(t)
		r.core.Events.Emit(event.OnDendroCore, t, a)
		r.core.Log.NewEvent(
			"dendro core spawned",
			glog.LogElementEvent,
			a.Info.ActorIndex,
		).
			Write("src", t.Src()).
			Write("expiry", r.core.F+t.Duration)
	}, dendrocore.Delay)
}
