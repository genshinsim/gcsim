package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillJumpFrames []int

func init() {
	skillJumpFrames = frames.InitAbilSlice(54) // skill
	skillJumpFrames[action.ActionAttack] = 4
	skillJumpFrames[action.ActionBurst] = 50
	skillJumpFrames[action.ActionDash] = 49
	skillJumpFrames[action.ActionJump] = 50
	skillJumpFrames[action.ActionWalk] = 47
	skillJumpFrames[action.ActionSwap] = 48
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		if c.Core.Player.LastAction.Type == action.ActionDash {
			c.reduceNightsoulPoints(14) // total 24, 10 from dash, 14 from dash jump
		} else {
			c.reduceNightsoulPoints(2)
		}

		switch c.Core.Player.AnimationHandler.CurrentState() {
		case action.DashState, action.JumpState, action.WalkState:
			// use the previous momentum gain tasks
		default:
			// queue a new momentum gain task
			c.momentumSrc = c.Core.F
			c.QueueCharTask(c.momentumStackGain(c.momentumSrc), momentumDelay)
		}

		return action.Info{
			Frames:          frames.NewAbilFunc(skillJumpFrames),
			AnimationLength: skillJumpFrames[action.InvalidAction],
			CanQueueAfter:   skillJumpFrames[action.ActionWalk],
			State:           action.JumpState,
		}, nil
	}
	return c.Character.Jump(p)
}
