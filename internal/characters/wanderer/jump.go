package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var EJumpFrames []int

func init() {
	EJumpFrames = frames.InitAbilSlice(46)
	EJumpFrames[action.ActionAttack] = 22
	EJumpFrames[action.ActionCharge] = 22
	EJumpFrames[action.ActionSkill] = 5
	EJumpFrames[action.ActionBurst] = 5
	EJumpFrames[action.ActionDash] = 5
	EJumpFrames[action.ActionWalk] = 45
}

func (c *char) Jump(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	if c.StatusIsActive(skillKey) {
		return c.WindfavoredJump(p)
	}

	ai := c.Character.Jump(p)
	ai.Frames = func(next action.Action) int { return delay + ai.Frames(next) }
	ai.AnimationLength = delay + ai.AnimationLength
	ai.CanQueueAfter = delay + ai.CanQueueAfter

	return ai
}

func (c *char) WindfavoredJump(p map[string]int) action.ActionInfo {
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(EJumpFrames),
		AnimationLength: EJumpFrames[action.ActionJump],
		CanQueueAfter:   EJumpFrames[action.ActionSkill], // earliest cancel
		State:           action.JumpState,
	}
}
