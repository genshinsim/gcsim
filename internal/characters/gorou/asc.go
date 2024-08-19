package gorou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After using Juuga: Forward Unto Victory, all nearby party members' DEF is increased by 25% for 12s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1Key, 720),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return c.a1Buff, true
			},
		})
	}
}

// Gorou receives the following DMG Bonuses to his attacks based on his DEF:
//
// - Inuzaka All-Round Defense: Skill DMG increased by 156% of DEF
func (c *char) a4Skill() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.TotalDef() * 1.56
}

// Gorou receives the following DMG Bonuses to his attacks based on his DEF:
//
// - Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
func (c *char) a4Burst() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.TotalDef() * 0.156
}
