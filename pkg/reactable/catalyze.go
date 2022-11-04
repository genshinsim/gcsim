package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) TryAggravate(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	if r.Durability[ModifierQuicken] < ZeroDur {
		return
	}

	r.core.Events.Emit(event.OnAggravate, r.self, a)

	//em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = combat.Aggravate
	a.Info.FlatDmg += 1.15 * r.calcCatalyzeDmg(a.Info, em)
}

func (r *Reactable) TrySpread(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	if r.Durability[ModifierQuicken] < ZeroDur {
		return
	}

	r.core.Events.Emit(event.OnSpread, r.self, a)

	//em isn't snapshot
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	a.Info.Catalyzed = true
	a.Info.CatalyzedType = combat.Spread
	a.Info.FlatDmg += 1.25 * r.calcCatalyzeDmg(a.Info, em)
}

func (r *Reactable) TryQuicken(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Dendro:
		if r.Durability[ModifierElectro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Electro, a.Info.Durability, 1)
	case attributes.Electro:
		if r.Durability[ModifierDendro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Dendro, a.Info.Durability, 1)
	default:
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	r.core.Events.Emit(event.OnQuicken, r.self, a)

	//attach quicken aura; special amount
	r.attachQuicken(consumed)

	if r.Durability[ModifierHydro] >= ZeroDur {
		r.core.Tasks.Add(func() {
			r.tryQuickenBloom(a)
		}, 0)
	}
}

func (r *Reactable) attachQuicken(dur combat.Durability) {
	r.attachOverlapRefreshDuration(ModifierQuicken, dur, 12*dur+360)
}
