package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4BuffKey = "shenhe-c4"

func (c *char) c4() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("shenhe-c4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
				return nil, false
			}
			if !c.StatusIsActive(c4BuffKey) {
				c.c4count = 0
				return nil, false
			}
			c.c4bonus[attributes.DmgP] += 0.05 * float64(c.c4count)
			c.c4count = 0
			c.DeleteStatus(c4BuffKey)
			return c.c4bonus, true
		},
	})
}
