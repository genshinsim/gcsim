package xilonen

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() == action.DashState && c.canUseNightsoul() {
		c.reduceNightsoulPoints(20)
	}
	return c.Character.Jump(p)
}
