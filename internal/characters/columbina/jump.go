package columbina

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(player.XianyunAirborneBuff) {
		return action.Info{}, errors.New("xianyun jump not implemented")
	}

	f := 48
	if c.Core.Player.CurrentState() == action.DashState {
		f = 52
	}

	return action.Info{
		Frames:          func(action.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           action.JumpState,
	}, nil
}
