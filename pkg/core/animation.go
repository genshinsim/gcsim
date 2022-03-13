package core

import "github.com/genshinsim/gcsim/pkg/coretype"

func (c *Core) SetState(state coretype.AnimationState, duration int) {
	c.Emit(coretype.OnStateChange, c.state, state)
	c.state = state
	c.stateExpiry = c.F + duration
}

func (c *Core) ClearState() {
	c.state = coretype.Idle
	c.stateExpiry = c.F - 1
}

func (c *Core) State() coretype.AnimationState {

	if c.stateExpiry > c.F {
		return c.state
	}
	return coretype.Idle
}
