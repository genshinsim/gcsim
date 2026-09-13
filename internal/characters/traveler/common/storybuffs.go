package common

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const TrueMoonParam = "true_moon_story_buff"

func TravelerStoryBuffs(c *character.CharWrapper, p info.CharacterProfile) {
	// TravelerStoryBuffs applies buffs based on completed story quests
	//
	// base_atk_buff
	// 		buffs from completing Chapter III: Act I (Sumeru Archon Quest)
	// 		+3 base atk
	// skirk_story_buff
	// 		buffs from completing Crystallina Chapter: Act I (Skirk's Story Quest)
	// 		+7 base atk, +15 EM, +50 base HP
	// true_moon_buff
	//		buffs from completing Song of the Welkin Moon: Act VIII - True Moon (Nod Krai Archon Quest)
	//		we assume traveler has all resonances done
	// All buffs default to enabled
	baseAtkBuff, okBaseAtkBuff := p.Params["base_atk_buff"]
	skirkBuff, okSkirkBuff := p.Params["skirk_story_buff"]
	trueMoonBuff, okTrueMoonBuff := p.Params[TrueMoonParam]
	if !okBaseAtkBuff {
		baseAtkBuff = 1
	}
	if !okSkirkBuff {
		skirkBuff = 1
	}
	if !okTrueMoonBuff {
		trueMoonBuff = 1
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

	if trueMoonBuff == 1 {
		m[attributes.CR] += 0.1
		m[attributes.DEFP] += 0.2
		m[attributes.ER] += 0.2
		m[attributes.EM] += 60
		m[attributes.HPP] += 0.2
		m[attributes.ATKP] += 0.2
		m[attributes.CD] += 0.2
	}

	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase("traveler-story-quest-buffs", -1),
		Amount: func() []float64 {
			return m
		},
	})
}

// the CA buffs must be added in the specific traveler
