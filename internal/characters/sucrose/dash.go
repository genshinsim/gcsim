package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var dashFrames []int

// sucrose's dash can be cancelled by her E and Q, so we override it here. wtf sucrose
func (c *char) Dash(p map[string]int) action.ActionInfo {
	// call default implementation to handle stamina
	c.Character.Dash(p)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(dashFrames),
		AnimationLength: dashFrames[action.InvalidAction],
		CanQueueAfter:   dashFrames[action.ActionDash], // earliest cancel
		Post:            dashFrames[action.ActionDash], // earliest cancel
		State:           action.DashState,
	}
}
