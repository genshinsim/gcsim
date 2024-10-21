package xilonen

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		if c.Core.Player.LastAction.Type == action.ActionDash {
			c.reduceNightsoulPoints(20)
		}
	}
	return c.Character.Jump(p)
}
