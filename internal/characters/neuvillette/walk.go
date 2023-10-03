package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.chargeEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack: Equitable Judgement with Walk", c.CharWrapper.Base.Key)
	}
	return c.Character.Walk(p)
}
