package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) action.ActionInfo {
	// chaining dashes resets savedNormalCounter
	if c.Core.Player.CurrentState() == action.DashState {
		c.savedNormalCounter = 0
	}

	return c.Character.Dash(p)
}
