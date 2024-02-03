package keqing

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

func (c *char) HighPlungeAttack(map[string]int) (action.Info, error) {
	if c.Core.Player.Airborne() == player.AirborneXianyun {
		//TODO: implementation goes here
		return action.Info{}, errors.New("keqing: action high_plunge during Xianyun jump not yet implemented")
	}
	return action.Info{}, fmt.Errorf("%v: action high_plunge not implemented", c.CharWrapper.Base.Key)
}

func (c *char) LowPlungeAttack(map[string]int) (action.Info, error) {
	if c.Core.Player.Airborne() == player.AirborneXianyun {
		//TODO: implementation goes here
		return action.Info{}, errors.New("keqing: action low_plunge during Xianyun jump not yet implemented")
	}
	return action.Info{}, fmt.Errorf("%v: action low_plunge not implemented", c.CharWrapper.Base.Key)
}
