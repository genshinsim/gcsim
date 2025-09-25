package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (r *Reactable) TryAggravate(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}

	if r.Durability[info.ReactionModKeyQuicken] < info.ZeroDur {
		return false
	}

	r.core.Events.Emit(event.OnAggravate, r.self, a)

	// em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = info.ReactionTypeAggravate
	a.Info.FlatDmg += 1.15 * r.calcCatalyzeDmg(a.Info, em)
	return true
}

func (r *Reactable) TrySpread(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}

	if r.Durability[info.ReactionModKeyQuicken] < info.ZeroDur {
		return false
	}

	r.core.Events.Emit(event.OnSpread, r.self, a)

	// em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = info.ReactionTypeSpread
	a.Info.FlatDmg += 1.25 * r.calcCatalyzeDmg(a.Info, em)
	return true
}

func (r *Reactable) TryQuicken(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}

	var consumed info.Durability
	switch a.Info.Element {
	case attributes.Dendro:
		if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur {
			return false
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	case attributes.Electro:
		if r.Durability[info.ReactionModKeyDendro] < info.ZeroDur {
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

	if r.Durability[info.ReactionModKeyHydro] >= info.ZeroDur {
		r.core.Tasks.Add(func() {
			r.tryQuickenBloom(a)
		}, 0)
	}

	return true
}

func (r *Reactable) attachQuicken(dur info.Durability) {
	r.attachOverlapRefreshDuration(info.ReactionModKeyQuicken, dur, 12*dur+360)
}
