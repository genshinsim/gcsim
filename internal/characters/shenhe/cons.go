package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("shenhe-c4", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return nil, false
		}
		if c.Core.F >= c.c4expiry {
			return nil, false
		}
		m[attributes.DmgP] += 0.05 * float64(c.c4count)
		c.c4count = 0
		c.c4expiry = 0
		return m, true
	})
}
