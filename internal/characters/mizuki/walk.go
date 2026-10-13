package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(dreamDrifterStateKey) {
		travel, _ := p["travel"]
		c.cloudTravel = travel
		c.cloudTravel = max(c.cloudTravel, 0)
		c.cloudTravel = min(c.cloudTravel, 120)
		skipTicks, _ := p["skip_ticks"]
		c.cloudSkipTicks = skipTicks
	}
	return c.Character.Walk(p)
}
