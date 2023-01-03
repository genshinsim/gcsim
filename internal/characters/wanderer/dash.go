package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var (
	a4Release   = []int{16, 18, 21, 25}
	dashFramesE []int
)

const a4Hitmark = 30

func init() {
	dashFramesE = frames.InitAbilSlice(24)
	dashFramesE[action.ActionAttack] = 21
	dashFramesE[action.ActionCharge] = 21
	dashFramesE[action.ActionSkill] = 5
	dashFramesE[action.ActionDash] = 22
	dashFramesE[action.ActionJump] = 22
	dashFramesE[action.ActionWalk] = 22
}

func (c *char) Dash(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	if c.StatusIsActive(skillKey) {
		return c.WindfavoredDash(p)
	}

	ai := c.Character.Dash(p)
	ai.Frames = func(action action.Action) int { return delay + ai.Frames(action) }
	ai.AnimationLength = delay + ai.AnimationLength
	ai.CanQueueAfter = delay + ai.CanQueueAfter

	return ai
}

func (c *char) WindfavoredDash(p map[string]int) action.ActionInfo {
	ai := action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFramesE[next] },
		AnimationLength: dashFramesE[action.InvalidAction],
		CanQueueAfter:   dashFramesE[action.ActionSkill],
		State:           action.DashState,
	}

	if c.StatusIsActive(a4Key) {
		c.a4()
	} else {
		c.skydwellerPoints -= 15
	}

	return ai
}
