package iansan

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() && c.Core.Player.CurrentState() == action.DashState {
		c.reduceNightsoulPoints(6)
	}
	return c.Character.Jump(p)
}
