package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryMelt(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Pyro:
		if r.Durability[info.ReactionModKeyCryo] < info.ZeroDur && r.Durability[info.ReactionModKeyFrozen] < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 2)
		f := r.reduce(attributes.Frozen, a.Info.Durability, 2)
		if f > consumed {
			consumed = f
		}
		a.Info.AmpMult = 2.0
	case attributes.Cryo:
		if r.Durability[info.ReactionModKeyPyro] < info.ZeroDur && r.Durability[info.ReactionModKeyBurning] < info.ZeroDur {
			return false
		}
		r.reduce(attributes.Pyro, a.Info.Durability, 0.5)
		a.Info.AmpMult = 1.5
		r.burningCheck()
	default:
		// should be here
		return false
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	a.Info.Amped = true
	a.Info.AmpType = info.ReactionTypeMelt
	r.core.Events.Emit(event.OnMelt, r.self, a)
	return true
}
