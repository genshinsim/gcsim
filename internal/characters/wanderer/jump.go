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

func (c *char) Jump(p map[string]int) (action.Info, error) {
	delay := c.checkForSkillEnd()

	if c.StatusIsActive(SkillKey) {
		return c.WindfavoredJump(p)
	}

	ai, err := c.Character.Jump(p)
	ai.Frames = func(next action.Action) int { return delay + ai.AnimationLength } // jump has static frames
	ai.AnimationLength = delay + ai.AnimationLength
	ai.CanQueueAfter = delay + ai.CanQueueAfter

	return ai, err
}

func (c *char) WindfavoredJump(p map[string]int) (action.Info, error) {
	return action.Info{
		Frames:          frames.NewAbilFunc(EJumpFrames),
		AnimationLength: EJumpFrames[action.ActionJump],
		CanQueueAfter:   EJumpFrames[action.ActionSkill], // earliest cancel
		State:           action.JumpState,
	}, nil
}
