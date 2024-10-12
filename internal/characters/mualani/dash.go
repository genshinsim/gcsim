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
