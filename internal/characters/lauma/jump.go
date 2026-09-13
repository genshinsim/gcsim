package lauma

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(deerStatusKey) {
		return action.Info{}, errors.New("jump in lauma spirit envoy state isn't implemented")
	}

	return c.Character.Jump(p)
}
