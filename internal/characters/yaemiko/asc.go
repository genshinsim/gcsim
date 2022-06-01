package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("yaemiko-a1", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		// only trigger on elemental art damage
		if atk.Info.AttackTag != combat.AttackTagElementalArt {
			return nil, false
		}
		m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
		return m, true
	})
}
