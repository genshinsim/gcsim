package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var skillFrames []int

const skillStart = 11

func init() {
	skillFrames = frames.InitAbilSlice(34)
	skillFrames[action.ActionAttack] = 22
	skillFrames[action.ActionAim] = 22 // uses attack frames
	skillFrames[action.ActionBurst] = 23
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 31
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	c.Core.Status.Add("yoimiyaskill", 600+skillStart) //activate for 10
	if c.Core.Status.Duration("yoimiyaa1") == 0 {
		c.a1stack = 0
	}

	c.SetCDWithDelay(action.ActionSkill, 1080, 11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) onExit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == c.Index && next != c.Index {
			if c.Core.Status.Duration("yoimiyaskill") > 0 {
				c.Core.Status.Delete("yoimiyaskill")
			}
		}
		return false
	}, "yoimiya-exit")
}
