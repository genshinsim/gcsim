package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (r *Reactable) TryAggravate(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	if r.Durability[model.Element_Overdose] < ZeroDur {
		return false
	}

	r.core.Events.Emit(event.OnAggravate, r.self, a)

	// em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = reactions.Aggravate
	a.Info.FlatDmg += 1.15 * r.calcCatalyzeDmg(a.Info, em)
	return true
}

func (r *Reactable) TrySpread(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	if r.Durability[model.Element_Overdose] < ZeroDur {
		return false
	}

	r.core.Events.Emit(event.OnSpread, r.self, a)

	// em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = reactions.Spread
	a.Info.FlatDmg += 1.25 * r.calcCatalyzeDmg(a.Info, em)
	return true
}

func (r *Reactable) TryQuicken(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	var consumed reactions.Durability
	switch a.Info.Element {
	case attributes.Dendro:
		if r.Durability[model.Element_Electric] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	case attributes.Electro:
		if r.Durability[model.Element_Grass] < ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Dendro, a.Info.Durability, 1)
	default:
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.core.Events.Emit(event.OnQuicken, r.self, a)

	// attach quicken aura; special amount
	r.attachQuicken(consumed)

	if r.Durability[model.Element_Water] >= ZeroDur {
		r.core.Tasks.Add(func() {
			r.tryQuickenBloom(a)
		}, 0)
	}

	return true
}

func (r *Reactable) attachQuicken(dur reactions.Durability) {
	r.attachOverlapRefreshDuration(model.Element_Overdose, dur, 12*dur+360)
}
