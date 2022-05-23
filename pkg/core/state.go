package core

type AnimationState int

const (
	Idle AnimationState = iota
	NormalAttackState
	ChargeAttackState
	PlungeAttackState
	SkillState
	BurstState
	AimState
	DashState
	JumpState
)

func (c *Core) SetState(state AnimationState, duration int) {
	c.Events.Emit(OnStateChange, c.state, state)
	c.state = state
	c.stateExpiry = c.F + duration
}

func (c *Core) ClearState() {
	c.state = Idle
	c.stateExpiry = c.F - 1
}

func (c *Core) State() AnimationState {

	if c.stateExpiry > c.F {
		return c.state
	}

	return Idle
}
