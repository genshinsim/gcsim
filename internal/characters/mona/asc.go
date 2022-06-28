package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

//Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{Base: modifier.NewBase("mona-a4", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		m[attributes.HydroP] = .2 * atk.Snapshot.Stats[attributes.ER]
		return m, true
	}})
}
