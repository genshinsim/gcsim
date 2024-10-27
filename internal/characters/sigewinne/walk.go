package sigewinne

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.burstEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Walk", c.CharWrapper.Base.Key)
	}
	return c.Character.Walk(p)
}
