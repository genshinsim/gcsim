package sayu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// Yoohoo Art: Fuuin Dash gains the following effects:
// DMG of Fuufuu Whirlwind Kick in Press Mode increased by 3.3%.
// Every 0.5s in the Fuufuu Windwheel state will increase the DMG of this Fuufuu
// Whirlwind Kick by 3.3%. The maximum DMG increase possible through this method
// is 66%.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("sayu-c2", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.ActorIndex != c.Index {
			return nil, false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return nil, false
		}
		m[attributes.DmgP] = c.c2Bonus
		//reset bonus back to 0
		c.c2Bonus = 0
		return m, true
	})
}
