package sara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// A1 is implemented in aimed.go:
// While in the Crowfeather Cover state provided by Tengu Stormcall, Aimed Shot charge times are decreased by 60%.

// not strictly required but in case in future we implement player getting hit
const a4ICDKey = "sara-a4-icd"

// When Tengu Juurai: Ambush hits opponents, Kujou Sara will restore 1.2 Energy to all party members for every 100% Energy Recharge she has. This effect can be triggered once every 3s.
//
// - according to library finding, text description is inaccurate
//
// - it's more like for every 1% of ER, she grants 0.012 flat energy
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(a4ICDKey) {
			return
		}
		c.AddStatus(a4ICDKey, 180, true)

		energyAddAmt := 1.2 * c.NonExtraStat(attributes.ER)
		for _, char := range c.Core.Player.Chars() {
			char.AddEnergy("sara-a4", energyAddAmt)
		}
	}
}
