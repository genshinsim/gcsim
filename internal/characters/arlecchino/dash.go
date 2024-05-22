package arlecchino

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.swapError {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack with Swap", c.CharWrapper.Base.Key)
	}

	c.chargeEarlyCancelled = false
	return c.Character.Dash(p)
}
