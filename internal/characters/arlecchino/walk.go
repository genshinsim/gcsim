package arlecchino

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.swapError {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack with Swap", c.CharWrapper.Base.Key)
	}

	if c.chargeEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack with Walk", c.CharWrapper.Base.Key)
	}
	return c.Character.Walk(p)
}
