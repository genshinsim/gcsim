package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var dashFrames []int

func init() {
	dashFrames = frames.InitAbilSlice(21)
	dashFrames[action.ActionSwap] = 20
	dashFrames[action.ActionSkill] = 2
	dashFrames[action.ActionBurst] = 1
	dashFrames[action.ActionCharge] = 20
}

func (c *char) Dash(p map[string]int) action.Info {
	// call default implementation to handle stamina
	c.Character.Dash(p)
	return action.Info{
		Frames:          frames.NewAbilFunc(dashFrames),
		AnimationLength: dashFrames[action.InvalidAction],
		CanQueueAfter:   dashFrames[action.ActionBurst],
		State:           action.DashState,
	}
}
