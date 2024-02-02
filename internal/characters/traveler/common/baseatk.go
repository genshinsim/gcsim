package common

import "github.com/genshinsim/gcsim/pkg/core/info"

// doing Chapter III: Act I of sumeru archon quest buffs base atk by 3
func TravelerBaseAtkIncrease(p info.CharacterProfile) float64 {
	baseAtkBuff, ok := p.Params["base_atk_buff"]
	if !ok {
		baseAtkBuff = 1
	}
	if baseAtkBuff == 1 {
		return 3
	}
	return 0
}
