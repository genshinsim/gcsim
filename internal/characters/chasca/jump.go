package chasca

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return action.Info{}, errors.New("chasca jump in nightsoul blessing not implemented")
	}
	return c.Character.Jump(p)
}
