package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (r *Reactable) tryQuicken(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	//quicken durability should be = amount consumed = min(existing, applying)
	//
	switch a.Info.Element {
	case attributes.Dendro:
		if r.Durability[attributes.Electro] < ZeroDur {
			return
		}
	case attributes.Electro:
		if r.Durability[attributes.Dendro] < ZeroDur {
			return
		}
	}

}
