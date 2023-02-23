package travelergeo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// Reduces Starfell Sword's CD by 2s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.skillCD -= 2 * 60
}

// The final hit of a Normal Attack combo triggers a collapse, dealing 60% of ATK as AoE Geo DMG.
func (c *char) a4() {
	if c.Base.Ascension < 4 || c.NormalCounter != c.NormalHitNum-1 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Frenzied Rockslide (A4)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       0.6,
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.2}, 2.4),
			0,
			0,
		)
	}, a4Hitmark[c.gender])
}
