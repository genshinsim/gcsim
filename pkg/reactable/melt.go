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
	if r.Durability[attributes.Frozen] > ZeroDur {
		return
	}
	switch a.Info.Element {
	case attributes.Pyro:
		if r.Durability[attributes.Cryo] < ZeroDur {
			return
		}
		r.reduce(attributes.Cryo, a.Info.Durability, 2)
		a.Info.AmpMult = 2.0
	case attributes.Cryo:
		if r.Durability[attributes.Pyro] < ZeroDur {
			return
		}
		r.reduce(attributes.Pyro, a.Info.Durability, 0.5)
		a.Info.AmpMult = 1.5
	default:
		//should be here
		return
	}
	//there shouldn't be anything else to react with if not frozen
	a.Info.Durability = 0
	a.Info.Amped = true
	a.Info.AmpType = combat.Melt
	r.core.Events.Emit(event.OnMelt, r.self, a)
}

func (r *Reactable) tryMeltFrozen(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[attributes.Frozen] < ZeroDur {
		return
	}
	switch a.Info.Element {
	case attributes.Pyro:
		//TODO: the assumption here is we first reduce cryo, and if there's any
		//src durability left, we reduce frozen. note that it's still only one
		//melt reaction
		a.Info.Durability -= r.reduce(attributes.Cryo, a.Info.Durability, 2)
		r.reduce(attributes.Frozen, a.Info.Durability, 2)
		a.Info.AmpMult = 2.0
	default:
		//should be here
		return
	}
	//durability not wiped out because we can potentially react with hydro still
	a.Info.Amped = true
	a.Info.AmpType = combat.Melt
	r.core.Events.Emit(event.OnMelt, r.self, a)
}
