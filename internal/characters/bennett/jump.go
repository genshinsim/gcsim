package bennett

import (
	"github.com/genshinsim/gcsim/internal/characters/xianyun"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var xianyunJumpFrames []int

func init() {
	xianyunJumpFrames = frames.InitAbilSlice(58)
	xianyunJumpFrames[action.ActionHighPlunge] = 6
	xianyunJumpFrames[action.ActionLowPlunge] = 5
}

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(xianyun.StarwickerKey) {
		c.Core.Player.SetAirborne(player.AirborneXianyun)
		return action.Info{
			Frames:          frames.NewAbilFunc(xianyunJumpFrames),
			AnimationLength: xianyunJumpFrames[action.InvalidAction],
			CanQueueAfter:   xianyunJumpFrames[action.ActionLowPlunge], // earliest cancel
			State:           action.JumpState,
		}, nil
	}
	return c.Character.Jump(p)
}
