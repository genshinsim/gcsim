package mualani

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Walk(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		switch c.Core.Player.AnimationHandler.CurrentState() {
		case action.DashState, action.JumpState, action.WalkState:
			// use the previous momentum gain tasks
		default:
			// queue a new momentum gain task
			c.momentumSrc = c.Core.F
			c.QueueCharTask(c.momentumStackGain(c.momentumSrc), momentumDelay)
		}
	}
	return c.Character.Walk(p)
}
