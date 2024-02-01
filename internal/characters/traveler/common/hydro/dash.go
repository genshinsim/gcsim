package hydro

import "github.com/genshinsim/gcsim/pkg/core/action"

func (c *Traveler) Dash(p map[string]int) (action.Info, error) {
	if c.Base.Ascension >= 1 {
		count := 0
		if p["pickup_droplets"] > 0 {
			count = p["pickup_droplets"]
		}
		c.a1PickUp(count)
	}

	return c.Character.Dash(p)
}
