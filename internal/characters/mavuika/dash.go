package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

// Special bike-dash will be here
func (c *char) Dash(p map[string]int) (action.Info, error) {
	return c.Character.Dash(p)
}
