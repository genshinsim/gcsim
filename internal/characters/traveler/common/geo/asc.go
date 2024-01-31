package geo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// Reduces Starfell Sword's CD by 2s.
func (c *Traveler) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.skillCD -= 2 * 60
}

// The final hit of a Normal Attack combo triggers a collapse, dealing 60% of ATK as AoE Geo DMG.
func (c *Traveler) a4() {
	if c.Base.Ascension < 4 || c.NormalCounter != c.NormalHitNum-1 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Frenzied Rockslide (A4)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   13.5,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       0.6,
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1.2}, 2.4),
			0,
			0,
		)
	}, a4Hitmark[c.gender])
}
