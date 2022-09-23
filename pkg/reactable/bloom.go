package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

func (r *Reactable) tryBloom(a *combat.AttackEvent) {
	//can be hydro bloom, dendro bloom, or quicken bloom
	if a.Info.Durability < ZeroDur {
		return
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		//this part is annoying. bloom will happen if any of the dendro like aura is present
		//so we gotta check for all 3...
		switch {
		case r.Durability[ModifierDendro] > ZeroDur:
		case r.Durability[ModifierQuicken] > ZeroDur:
		case r.Durability[ModifierBurningFuel] > ZeroDur:
		default:
			return
		}
		//reduce only check for one element so have to call twice to check for quicken as well
		consumed = r.reduce(attributes.Dendro, a.Info.Durability, 0.5)
		f := r.reduce(attributes.Quicken, a.Info.Durability, 0.5)
		if f > consumed {
			consumed = f
		}
	case attributes.Dendro:
		if r.Durability[ModifierHydro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Hydro, a.Info.Durability, 2)
	case attributes.Quicken:
		//TODO: ?? how to handle this??
		if r.Durability[ModifierHydro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Quicken, a.Info.Durability, 2)
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.core.Events.Emit(event.OnBloom, r.self, a)
}

func (r *Reactable) tryHyperbloom(a *combat.AttackEvent) {

}

func (r *Reactable) tryBurgeon(a *combat.AttackEvent) {

}

type dendroCore struct {
	*gadget.Gadget
}
