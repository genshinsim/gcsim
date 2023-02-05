package layla

import (
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// While the Curtain of Slumber is active, the Deep Sleep effect will activate each time the Curtain gains 1 Night Star:
//
// - The Shield Strength of a character under the effect of the Curtain of Slumber increases by 6%.
//
// - This effect can have a maximum of 4 stacks and persists until the Curtain of Slumber disappears.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Player.Shields.AddShieldBonusMod("layla-a1", -1, func() (float64, bool) {
		if exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill); exist == nil {
			return 0, false
		}
		return float64(c.a1Stack) * 0.06, false
	})
}

// The DMG dealt by the Shooting Stars fired by Nights of Formal Focus is increased by 1.5% of Layla's Max HP.
func (c *char) a4() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return 0.015 * c.MaxHP()
}
