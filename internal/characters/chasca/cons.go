package chasca

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

func (c *char) c1() float64 {
	if c.Base.Cons < 1 {
		return 0.0
	}
	return 0.333
}
func (c *char) c1Conversion() attributes.Element {
	if c.Base.Cons < 1 {
		return attributes.Anemo
	}
	return c.a1Conversion()
}
func (c *char) c2A1Stack() int {
	if c.Base.Cons < 2 {
		return 0
	}
	return 1
}

func (c *char) c2cb(src int) combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(ac combat.AttackCB) {
		if c.c2Src == src {
			return
		}
		c.c2Src = src

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Shining Shadowhunt Shell (C2)",
			AttackTag:      attacks.AttackTagExtra,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        ac.AttackEvent.Info.Element,
			Durability:     25,
			Mult:           4,
		}
		// TODO: I can't find C2 AoE info
		ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -3}, 8, 120)
		c.Core.QueueAttack(ai, ap, 0, 1)
	}
}

func (c *char) c4cb(src int) combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(ac combat.AttackCB) {
		c.AddEnergy("chasca-c4", 1.5)
		if c.c4Src == src {
			return
		}
		c.c4Src = src

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Radiant Shadowhunt Shell (C2)",
			AttackTag:      attacks.AttackTagExtra,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        ac.AttackEvent.Info.Element,
			Durability:     25,
			Mult:           4,
		}
		// TODO: I can't find C4 AoE info
		ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -3}, 8, 120)
		c.Core.QueueAttack(ai, ap, 0, 1)
	}
}
