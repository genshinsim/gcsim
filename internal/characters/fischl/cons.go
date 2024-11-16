package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

func (c *char) c6Wave() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fischl C6",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupFischl,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       0.3,
	}

	// C6 uses Oz Snapshot
	c.Core.QueueAttackWithSnap(
		ai,
		c.ozSnapshot.Snapshot,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -1},
			0.1,
			1,
		),
		c.ozTravel,
	)
}
