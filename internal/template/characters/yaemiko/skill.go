package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillFrames []int

// kitsune spawn frame
const skillStart = 34

func (c *char) Skill(p map[string]int) action.ActionInfo {

	c.Core.Tasks.Add(func() { c.makeKitsune() }, skillStart)
	c.SetCD(action.ActionSkill, 4*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
