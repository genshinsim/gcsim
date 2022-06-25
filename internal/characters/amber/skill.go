package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillFrames []int

const skillStart = 35

func init() {
	skillFrames = frames.InitAbilSlice(35)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]

	c.Core.Tasks.Add(func() {
		c.makeBunny()
	}, skillStart+hold)

	if c.Base.Cons >= 4 {
		c.SetCD(action.ActionSkill, 720)
	} else {
		c.SetCD(action.ActionSkill, 900)
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillFrames[next] + hold },
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillStart + hold,
		State:           action.SkillState,
	}
}
