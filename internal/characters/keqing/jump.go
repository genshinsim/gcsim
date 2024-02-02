package keqing

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	if c.Core.Status.Duration(player.XianyunAirborneBuff) > 0 {
		// TODO: implementation to set airborne flag
		c.Core.Player.SetAirborne(player.AirborneXianyun)
		return action.Info{}, errors.New("keqing buffed jump not yet implemented")
	}
	return c.Character.Jump(p)
}
