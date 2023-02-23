package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

func (r *Reactable) TryVaporize(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Pyro:
		//make sure there's hydro
		if r.Durability[ModifierHydro] < ZeroDur {
			return false
		}
		//if there's still frozen left don't try to vape
		//game actively rejects vaporize reaction if frozen is present
		if r.Durability[ModifierFrozen] > ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Hydro, a.Info.Durability, .5)
		a.Info.AmpMult = 1.5
	case attributes.Hydro:
		//make sure there's pyro to vape; no coexistance with pyro (yet)
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Pyro, a.Info.Durability, 2)
		a.Info.AmpMult = 2
		r.burningCheck()
	default:
		//should be here
		return false
	}
	//there shouldn't be anything else to react with
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	a.Info.Amped = true
	a.Info.AmpType = reactions.Vaporize
	r.core.Events.Emit(event.OnVaporize, r.self, a)
	return true
}
