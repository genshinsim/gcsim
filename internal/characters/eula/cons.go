package eula

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25
	c.AddAttackMod("eula-c4", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.Abil != "Glacial Illumination (Lightfall)" {
			return nil, false
		}
		if !c.Core.Combat.DamageMode {
			return nil, false
		}
		if t.HP()/t.MaxHP() >= 0.5 {
			return nil, false
		}
		return m, true
	})
}
