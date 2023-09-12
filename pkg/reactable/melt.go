package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

func (r *Reactable) TryMelt(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed reactions.Durability
	switch a.Info.Element {
	case attributes.Pyro:
		if r.Durability[ModifierCryo] < ZeroDur && r.Durability[ModifierFrozen] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 2)
		f := r.reduce(attributes.Frozen, a.Info.Durability, 2)
		if f > consumed {
			consumed = f
		}
		a.Info.AmpMult = 2.0
	case attributes.Cryo:
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
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
	a.Info.AmpType = reactions.Melt
	r.core.Events.Emit(event.OnMelt, r.self, a)
	return true
}
