package alyosha

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

// Radiance: Stellar-Conduct: The Hunter's Precision effect obtained when the Hunter's Mark is
// activated will also increase Stellar-Conduct reaction DMG dealt by the character when on the
// field by 20%.
func (c *char) stellarInit() {
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("alyosha-ssc", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.isStellarRadiance() {
					return 0
				}
				if !c.StatusIsActive(skillBuffKey) {
					return 0
				}
				if ai.AttackTag != attacks.AttackTagDirectStellarConduct {
					return 0
				}
				if c.Core.Player.Active() != char.Index() {
					return 0
				}
				return 0.2 * float64(c.skillStacks)
			},
		})
	}
}

func (c *char) isStellarRadiance() bool {
	return c.StatusIsActive(reactable.PolestarFieldKey)
}
