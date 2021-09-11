package core

type AnimationState int

const (
	Idle AnimationState = iota
	Normal1State
)

func (c *Core) SetState(state AnimationState, duration int) {
	c.state = state
	c.stateExpiry = c.F + duration
}

func (c *Core) State() AnimationState {

	if c.stateExpiry > c.F {
		return c.state
	}

	return Idle
}
