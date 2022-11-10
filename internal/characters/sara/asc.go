package sara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

//not strictly required but in case in future we implement player getting hit
const a4ICDKey = "sara-a4-icd"

// Implements A4 energy regen.
// According to library finding, text description is inaccurate
// it's more like for every 1% of ER, she grants 0.012 flat energy
func (c *char) a4(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(a4ICDKey) {
		return
	}
	c.AddStatus(a4ICDKey, 180, true)
	energyAddAmt := 1.2 * (1 + c.Stat(attributes.ER))
	for _, char := range c.Core.Player.Chars() {
		char.AddEnergy("sara-a4", energyAddAmt)
	}
}
