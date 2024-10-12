package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillDashFrames []int

func init() {
	skillDashFrames = frames.InitAbilSlice(24) // dash
	skillDashFrames[action.ActionAttack] = 3
	skillDashFrames[action.ActionSkill] = 2
	skillDashFrames[action.ActionBurst] = 4
	skillDashFrames[action.ActionJump] = 2
	skillDashFrames[action.ActionWalk] = 3
	skillDashFrames[action.ActionSwap] = 1
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.reduceNightsoulPoints(10)

		switch c.Core.Player.AnimationHandler.CurrentState() {
		case action.DashState, action.JumpState, action.WalkState:
			// use the previous momentum gain tasks
		default:
			// queue a new momentum gain task
			c.momentumSrc = c.Core.F
			c.QueueCharTask(c.momentumStackGain(c.momentumSrc), momentumDelay)
		}

		// assuming doesn't contribute to dash CD
		return action.Info{
			Frames:          frames.NewAbilFunc(skillDashFrames),
			AnimationLength: skillDashFrames[action.InvalidAction],
			CanQueueAfter:   skillDashFrames[action.ActionSwap],
			State:           action.DashState,
		}, nil
	}
	return c.Character.Dash(p)
}
