package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	// anything but NA/E -> E should reset savedNormalCounter
	switch c.Core.Player.LastAction.Type {
	case action.ActionAttack:
	case action.ActionSkill:
	default:
		c.savedNormalCounter = 0
	}

	return c.Character.Dash(p)
}
