package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var StellarconductMult = []float64{1, 1.45, 1.5, 1.55, 1.6, 1.65, 1.7, 1.75, 1.8, 1.85, 1.9, 1.95, 2}

const (
	PolestarFieldKey       = "polestar-field"
	StellarConductShredKey = PolestarFieldKey + "-phys-shred"
)

func (r *Reactable) TryStellarConduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for non frozen one
	if r.GetAuraDurability(info.ReactionModKeyFrozen) >= info.ZeroDur {
		return false
	}
	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Electro:
		if r.GetAuraDurability(info.ReactionModKeyCryo) < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Cryo, a.Info.Durability, 1)
	case attributes.Cryo:
		// could be ec potentially
		if r.GetAuraDurability(info.ReactionModKeyElectro) < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	default:
		return false
	}

	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.queueStellarConduct(a)
	return true
}

func (r *Reactable) TryFrozenStellarConduct(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// this is for frozen
	if r.GetAuraDurability(info.ReactionModKeyFrozen) < info.ZeroDur {
		return false
	}
	switch a.Info.Element {
	case attributes.Electro:
		// TODO: the assumption here is we first reduce cryo, and if there's any
		// src durability left, we reduce frozen. note that it's still only one
		// stellarconduct reaction
		a.Info.Durability -= r.reduce(attributes.Cryo, a.Info.Durability, 1)
		r.reduce(attributes.Frozen, a.Info.Durability, 1)
		a.Info.Durability = 0
		a.Reacted = true
	default:
		return false
	}

	r.queueStellarConduct(a)

	return false
}

func (r *Reactable) queueStellarConduct(a *info.AttackEvent) {
	r.core.Events.Emit(event.OnStellarConduct, r.self, a)
	var nearbyPolestarField *PolestarField
	for _, g := range r.core.Combat.Gadgets() {
		polestarField, ok := g.(*PolestarField)
		if !ok {
			continue
		}
		if !r.core.Combat.Player().IsWithinArea(polestarField.fieldArea) {
			continue
		}
		nearbyPolestarField = polestarField
		break
	}

	if nearbyPolestarField == nil {
		r.newPolestarField()
		return
	}

	nearbyPolestarField.resetDuration()
}
