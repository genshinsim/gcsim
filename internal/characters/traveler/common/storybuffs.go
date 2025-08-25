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
	storyBuff, ok := p.Params["traveler_story_buff"]
	if !ok {
		storyBuff = 2 // default to maximum buffs
	}

	switch storyBuff {
	case 0:
		// no buffs
		return
	case 1:
		// Chapter III: Act I buffs only
		m := make([]float64, attributes.EndStatType)
		m[attributes.BaseATK] = 3
		c.AddStatMod(character.StatMod{
			Base: modifier.NewBase("traveler-sumeru-quest-buff", -1),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	case 2:
		// Chapter III: Act I + Skirk's Story Quest buffs
		m := make([]float64, attributes.EndStatType)
		m[attributes.BaseATK] = 3 + 7 // 3 from Sumeru + 7 from Skirk
		m[attributes.EM] = 15         // 15 EM from Skirk
		m[attributes.BaseHP] = 50     // 50 base HP from Skirk
		c.AddStatMod(character.StatMod{
			Base: modifier.NewBase("traveler-story-quest-buffs", -1),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

// TravelerBaseAtkIncrease is deprecated, use TravelerStoryBuffs instead
// keeping for backward compatibility with old parameter name
func TravelerBaseAtkIncrease(c *character.CharWrapper, p info.CharacterProfile) {
	// Check if old parameter is being used
	if _, ok := p.Params["base_atk_buff"]; ok {
		baseAtkBuff := p.Params["base_atk_buff"]
		if baseAtkBuff == 1 {
			// Convert old behavior to new system
			newParams := make(map[string]int)
			for k, v := range p.Params {
				newParams[k] = v
			}
			newParams["traveler_story_buff"] = 1
			newProfile := p
			newProfile.Params = newParams
			TravelerStoryBuffs(c, newProfile)
		}
		return
	}

	// Use new system
	TravelerStoryBuffs(c, p)
}
