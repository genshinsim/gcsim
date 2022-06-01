package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

//Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("mona-a4", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		m[attributes.HydroP] = .2 * atk.Snapshot.Stats[attributes.ER]
		return m, true
	})
}
