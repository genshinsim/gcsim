package jahoda

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillDashFrames []int

func init() {
	skillDashFrames = frames.InitAbilSlice(10) // D -> N1. Frames needed
	skillDashFrames[action.ActionAim] = 10     // D -> Aim. Frames needed
	skillDashFrames[action.ActionBurst] = 10   // D -> Q. Frames needed
	skillDashFrames[action.ActionDash] = 10    // D -> D. Frames needed
	skillDashFrames[action.ActionJump] = 10    // D -> J. Frames needed
	skillDashFrames[action.ActionWalk] = 10    // D -> W. Frames needed
	skillDashFrames[action.ActionSwap] = 10    // D -> Swap. Frames needed
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionBurst], // earliest cancel, need checking
			State:           action.SkillState,
		}, nil
	}
	return c.Character.Dash(p)
}
