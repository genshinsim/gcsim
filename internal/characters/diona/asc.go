package diona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) a1() {
	c.Core.Player.AddStamPercentMod("diona-a1", -1, func(a action.Action) (float64, bool) {
		if c.Core.Player.Shields.Get(shield.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	})
}
