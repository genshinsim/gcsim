package neuvillette

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	c.chargeEarlyCancelled = false
	return c.Character.Dash(p)
}
