package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Decreases the number of opponents Sweeping Fervor must hit to trigger each level of shielding.
//
// - Shield Level 2: Lead-In requirement reduced to 1 opponent hit.
//
// - Shield Level 3: Rave requirement reduced to 2 opponents hit or more.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.shieldLevel2Requirement -= 1
	c.shieldLevel3Requirement -= 1
}

// Characters shielded by Sweeping Fervor deal 15% increased Physical DMG.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.PhyP] = 0.15
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("xinyan-a4", -1),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if !c.Core.Player.Shields.PlayerIsShielded() {
					return nil, false
				}
				shd := c.Core.Player.Shields.Get(shield.XinyanSkill)
				if shd == nil {
					return nil, false
				}
				return m, true
			},
		})
	}
}
