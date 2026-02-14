package lauma

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(deerStatusKey) {
		return action.Info{}, errors.New("cannot use dash in lauma spirit envoy state")
	}

	return c.Character.Dash(p)
}
