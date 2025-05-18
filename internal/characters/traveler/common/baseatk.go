package common

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// doing Chapter III: Act I of sumeru archon quest buffs base atk by 3
func TravelerBaseAtkIncrease(c *character.CharWrapper, p info.CharacterProfile) {
	baseAtkBuff, ok := p.Params["base_atk_buff"]
	if !ok {
		baseAtkBuff = 1
	}
	if baseAtkBuff != 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.BaseATK] = 3
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase("traveler-base-atk-buff", -1),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
