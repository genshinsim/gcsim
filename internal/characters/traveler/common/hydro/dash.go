package hydro

import "github.com/genshinsim/gcsim/pkg/core/action"

func (c *Traveler) Dash(p map[string]int) (action.Info, error) {
	if c.Base.Ascension >= 1 {
		count := max(p["pickup_droplets"], 0)
		c.a1PickUp(count)
	}

	return c.Character.Dash(p)
}
