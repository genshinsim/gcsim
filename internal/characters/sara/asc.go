package sara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// Implements A4 energy regen.
// According to library finding, text description is inaccurate
// it's more like for every 1% of ER, she grants 0.012 flat energy
func (c *char) a4(a combat.AttackCB) {
	if c.Core.F < c.a4LastProc {
		return
	}
	c.a4LastProc = c.Core.F + 180
	energyAddAmt := 1.2 * (1 + c.Stat(attributes.ER))
	for _, char := range c.Core.Player.Chars() {
		char.AddEnergy("sara-a4", energyAddAmt)
	}
}
