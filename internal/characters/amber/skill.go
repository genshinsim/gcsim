package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillFrames []int

const skillStart = 5 // cd start and bunny release frame on tap e
const bunnyLand = 45 // bunny land/spawn on tap e

func init() {
	skillFrames = frames.InitAbilSlice(33) // E -> E (C1 only)
	skillFrames[action.ActionAttack] = 32  // E -> N1
	skillFrames[action.ActionBurst] = 32   // E -> Q
	skillFrames[action.ActionDash] = 8     // E -> D
	skillFrames[action.ActionJump] = 8     // E -> J
	skillFrames[action.ActionSwap] = 23    // E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := p["hold"]

	c.Core.Tasks.Add(func() {
		c.makeBunny()
	}, bunnyLand+hold)

	if c.Base.Cons >= 4 {
		c.SetCDWithDelay(action.ActionSkill, 720, skillStart+hold)
	} else {
		c.SetCDWithDelay(action.ActionSkill, 900, skillStart+hold)
	}

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] + hold },
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillStart + hold,
		State:           action.SkillState,
	}, nil
}
