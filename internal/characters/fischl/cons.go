package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const c6HexereiKey = "fischl-c6-hexerei"

func (c *char) c6Wave() {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Evernight Raven (C6)",
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
			info.Point{Y: -1},
			0.1,
			1,
		),
		c.ozTravel,
	)

	if c.IsHexerei && c.Core.Player.GetHexereiCount() >= 2 {
		c.AddStatus(c6HexereiKey, 10*60, true)
	}
}

func (c *char) c6HexBonus() float64 {
	if c.Base.Cons < 6 {
		return 0
	}

	if !c.StatusIsActive(c6HexereiKey) {
		return 0
	}

	return 1
}
