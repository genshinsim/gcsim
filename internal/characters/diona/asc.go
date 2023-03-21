package diona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// Characters shielded by Icy Paws have their Movement SPD increased by 10% and their Stamina Consumption decreased by 10%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Player.AddStamPercentMod("diona-a1", -1, func(_ action.Action) (float64, bool) {
		if c.Core.Player.Shields.Get(shield.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	})
}

// A4 is not implemented:
// TODO: Opponents who enter the AoE of Signature Mix have 10% decreased ATK for 15s.
