package common

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// TravelerStoryBuffs applies buffs based on completed story quests
// 0 - no buffs
// 1 - buffs from completing "Chapter III: Act I of sumeru archon quest" (+3 base atk)
// 2 (default) - buff from 1 + buffs from completing "Skirk's Story Quest" (+7 additional base atk, +15 EM, +50 base HP)
func TravelerStoryBuffs(c *character.CharWrapper, p info.CharacterProfile) {
	baseAtkBuff, okBaseAtkBuff := p.Params["base_atk_buff"]
	skirkBuff, okSkirkBuff := p.Params["skirk_story_buff"]
	if !okBaseAtkBuff {
		baseAtkBuff = 1
	}
	if !okSkirkBuff {
		skirkBuff = 1 // default to maximum buffs
	}

	m := make([]float64, attributes.EndStatType)
	if baseAtkBuff == 1 {
		m[attributes.BaseATK] += 3
	}
	if skirkBuff == 1 {
		m[attributes.BaseATK] += 7
		m[attributes.EM] += 15
		m[attributes.BaseHP] += 50
	}
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase("traveler-story-quest-buffs", -1),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
