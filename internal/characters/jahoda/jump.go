package jahoda

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func init() {
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}, nil
	}
	return c.Character.Jump(p)
}
