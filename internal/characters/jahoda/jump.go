package jahoda

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillJumpFrames []int

func init() {
	skillJumpFrames = frames.InitAbilSlice(10) // J -> N1. Frames needed
	skillJumpFrames[action.ActionAim] = 10     // J -> Aim. Frames needed
	skillJumpFrames[action.ActionBurst] = 10   // J -> Q. Frames needed
	skillJumpFrames[action.ActionDash] = 10    // J -> D. Frames needed
	skillJumpFrames[action.ActionJump] = 10    // J -> J. Frames needed
	skillJumpFrames[action.ActionWalk] = 10    // J -> W. Frames needed
	skillJumpFrames[action.ActionSwap] = 10    // J -> Swap. Frames needed
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionBurst], // earliest cancel, need checking
			State:           action.SkillState,
		}, nil
	}
	return c.Character.Jump(p)
}
