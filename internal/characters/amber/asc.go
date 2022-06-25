package amber

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = .1
	c.AddAttackMod("amber-a1", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		return m, atk.Info.AttackTag == combat.AttackTagElementalBurst
	})
}

func (c *char) a4(a combat.AttackCB) {
	if !a.AttackEvent.Info.HitWeakPoint {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15
	c.AddStatMod("amber-a4", 600, attributes.ATKP, func() ([]float64, bool) {
		return m, true
	})
}
