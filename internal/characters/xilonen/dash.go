package xilonen

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.c6()

		if !c.StatusIsActive(c6key) {
			c.reduceNightsoulPoints(5.0 * c.c1ValMod())
		}
	}
	return c.Character.Dash(p)
}
