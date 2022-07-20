package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Characters shielded by Sweeping Fervor deal 15% increased Physical DMG.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PhyP] = 0.15
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("xinyan-a4", -1),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if !c.Core.Player.Shields.PlayerIsShielded() {
					return nil, false
				}
				shd := c.Core.Player.Shields.Get(shield.ShieldXinyanSkill)
				if shd == nil {
					return nil, false
				}
				return m, true
			},
		})
	}
}
