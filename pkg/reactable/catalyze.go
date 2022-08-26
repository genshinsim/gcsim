package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryAggravate(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	switch a.Info.Element {
	case attributes.Electro:
		if r.Durability[attributes.Quicken] < ZeroDur {
			return
		}
		// Aggravate doesn't consume any gauge
	default:
		//Trigger must be Electro. Maybe add a panic/error?
		return
	}

	//trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events.Emit(event.OnAggravate, r.self, a)

	//em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.FlatDmg += 1.15 * r.calcCatalyzeDmg(a.Info, em)
	a.Info.Amped = true
	a.Info.AmpType = combat.Aggravate
}

func (r *Reactable) trySpread(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	switch a.Info.Element {
	case attributes.Dendro:
		if r.Durability[attributes.Quicken] < ZeroDur {
			return
		}
		// Spread doesn't consume any gauge
	default:
		//Trigger must be Dendro. Maybe add a panic/error?
		return
	}

	//trigger event before attack is queued. this gives time for other actions to modify it
	r.core.Events.Emit(event.OnSpread, r.self, a)

	//em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.FlatDmg += 1.25 * r.calcCatalyzeDmg(a.Info, em)
	a.Info.Amped = true
	a.Info.AmpType = combat.Spread
}

func (r *Reactable) tryQuicken(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	switch a.Info.Element {
	case attributes.Dendro:
		//if electro exists we'll trigger quicken regardless if quicken already coexists
		if r.Durability[attributes.Electro] > ZeroDur {
			consumed := r.triggerQuicken(r.Durability[attributes.Electro], a.Info.Durability)
			r.Durability[attributes.Electro] = max(r.Durability[attributes.Electro]-consumed, 0)
			a.Info.Durability = max(a.Info.Durability-consumed, 0)
			r.core.Events.Emit(event.OnQuicken, r.self, a)
			return
		} else {
			//try refill first - this will use up all durability if ok
			r.tryRefill(a.Info.Element, &a.Info.Durability)
			//otherwise attach
			r.tryAttach(a.Info.Element, &a.Info.Durability)
		}
	case attributes.Electro:
		//if electro exists we'll trigger quicken regardless if quicken already coexists
		if r.Durability[attributes.Dendro] > ZeroDur {
			consumed := r.triggerQuicken(r.Durability[attributes.Dendro], a.Info.Durability)
			r.Durability[attributes.Dendro] = max(r.Durability[attributes.Dendro]-consumed, 0)
			a.Info.Durability = max(a.Info.Durability-consumed, 0)
			r.core.Events.Emit(event.OnQuicken, r.self, a)
			return
		} else {
			//try refill first - this will use up all durability if ok
			r.tryRefill(a.Info.Element, &a.Info.Durability)
			//otherwise attach
			r.tryAttach(a.Info.Element, &a.Info.Durability)
		}
	default:

		return
	}
}

//add to quicken durability and return amount of durability consumed
func (r *Reactable) triggerQuicken(a, b combat.Durability) combat.Durability {
	d := min(a, b)
	if d > r.Durability[attributes.Quicken] {
		r.Durability[attributes.Quicken] = d
		r.DecayRate[attributes.Quicken] = d / (12*d + 360.0)
	}
	return d
}
