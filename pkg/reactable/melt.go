package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryMelt(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Pyro:
		if r.Durability[ModifierCryo] < ZeroDur && r.Durability[ModifierFrozen] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 2)
		f := r.reduce(attributes.Frozen, a.Info.Durability, 2)
		if f > consumed {
			consumed = f
		}
		a.Info.AmpMult = 2.0
	case attributes.Cryo:
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
			return
		}
		r.reduce(attributes.Pyro, a.Info.Durability, 0.5)
		a.Info.AmpMult = 1.5
	default:
		//should be here
		return
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	a.Info.Amped = true
	a.Info.AmpType = combat.Melt
	r.core.Events.Emit(event.OnMelt, r.self, a)
}
