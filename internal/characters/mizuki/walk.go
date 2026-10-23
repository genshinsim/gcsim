package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(dreamDrifterStateKey) {
		c.cloudTravel = p["travel"]
		c.cloudTravel = max(c.cloudTravel, 0)
		c.cloudTravel = min(c.cloudTravel, 120)
		c.cloudSkipTicks = p["skip_ticks"]
	}
	return c.Character.Walk(p)
}
