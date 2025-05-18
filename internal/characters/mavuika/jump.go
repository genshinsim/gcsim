package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var bikeJumpFrames []int

func init() {
	bikeJumpFrames = frames.InitAbilSlice(42)
	bikeJumpFrames[action.ActionLowPlunge] = 16
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if !c.StatusIsActive(player.XianyunAirborneBuff) && c.armamentState == bike &&
		c.nightsoulState.HasBlessing() && c.Core.Player.CurrentState() == action.WalkState {
		c.canBikePlunge = true
		// Set plunge availability to false after jump duration
		// Plunge actions last long enough that we don't need a srcF check
		c.QueueCharTask(func() {
			c.canBikePlunge = false
		}, bikeJumpFrames[action.InvalidAction])
		return action.Info{
			Frames:          frames.NewAbilFunc(bikeJumpFrames),
			AnimationLength: bikeJumpFrames[action.InvalidAction],
			CanQueueAfter:   bikeJumpFrames[action.ActionLowPlunge], // earliest cancel
			State:           action.JumpState,
		}, nil
	}
	return c.Character.Jump(p)
}
