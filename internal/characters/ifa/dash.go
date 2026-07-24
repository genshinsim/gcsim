package ifa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillDashFrames []int

func init() {
	skillDashFrames = frames.InitAbilSlice(23)
	skillDashFrames[action.ActionAttack] = 17
	skillDashFrames[action.ActionBurst] = 15
	skillDashFrames[action.ActionSkill] = 9
	// TODO: fix skill hold frames
	// skillDashFrames[action.ActionSkillHoldFramesOnly] = 19
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.reduceNightsoulPoints(13.5)
		c.ApplyDashCD()

		return action.Info{
			Frames:          frames.NewAbilFunc(skillDashFrames),
			AnimationLength: skillDashFrames[action.InvalidAction],
			CanQueueAfter:   skillDashFrames[action.ActionSkill],
			State:           action.DashState,
		}, nil
	}
	return c.Character.Dash(p)
}
