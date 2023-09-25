package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var dashFrames []int

func init() {
	dashFrames = frames.InitAbilSlice(24)
	dashFrames[action.ActionSkill] = 1
	dashFrames[action.ActionBurst] = 1
}

// sucrose's dash can be cancelled by her E and Q, so we override it here. wtf sucrose
func (c *char) Dash(p map[string]int) (action.Info, error) {
	// call default implementation to handle stamina
	c.Character.Dash(p)
	return action.Info{
		Frames:          frames.NewAbilFunc(dashFrames),
		AnimationLength: dashFrames[action.InvalidAction],
		CanQueueAfter:   dashFrames[action.ActionSkill], // earliest cancel
		State:           action.DashState,
	}, nil
}
