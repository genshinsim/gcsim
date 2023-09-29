package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var burstJumpFrames []int

func init() {
	burstJumpFrames = frames.InitAbilSlice(58)
	burstJumpFrames[action.ActionHighPlunge] = 6
	burstJumpFrames[action.ActionLowPlunge] = 5
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(burstBuffKey) {
		return action.Info{
			Frames:          frames.NewAbilFunc(burstJumpFrames),
			AnimationLength: burstJumpFrames[action.InvalidAction],
			CanQueueAfter:   burstJumpFrames[action.ActionLowPlunge], // earliest cancel
			State:           action.JumpState,
		}, nil
	}
	return c.Character.Jump(p)
}
