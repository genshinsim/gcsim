package skirk

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
)

const skillDashLength = 11

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if !c.StatusIsActive(skillKey) {
		return c.Character.Dash(p)
	}

	// Execute dash CD logic
	c.ApplyDashCD()

	// consume stamina at end of the dash
	c.Core.Tasks.Add(func() {
		req := c.Core.Player.AbilStamCost(c.Index(), action.ActionDash, p)
		c.Core.Player.UseStam(req, action.ActionDash)
	}, skillDashLength)

	dashJumpLength := c.DashToJumpLength()
	return action.Info{
		Frames: func(a action.Action) int {
			switch a {
			case action.ActionJump:
				return dashJumpLength
			default:
				return skillDashLength
			}
		},
		AnimationLength: skillDashLength,
		CanQueueAfter:   dashJumpLength,
		State:           action.DashState,
	}, nil
}
