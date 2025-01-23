package chasca

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		// TODO: Figure out how to consume nightsoul here, and how this affects plunge timing
		return action.Info{}, errors.New("chasca jump in nightsoul blessing not implemented")
	}
	return c.Character.Jump(p)
}
