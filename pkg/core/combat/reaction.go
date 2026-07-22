package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type reactionBonusSrc interface {
	ReactBonus(atk info.AttackInfo) float64
}

func CalcReactionBaseDmg(lvl int) float64 {
	idx := lvl - 1
	idx = min(idx, 99)
	idx = max(idx, 0)
	return reactionLvlBase[idx]
}

func CalcLunarReactionDmg(lvl int, reactBonus float64, atk info.AttackInfo, em float64) float64 {
	var reactionMultiplier float64
	switch atk.AttackTag {
	case attacks.AttackTagReactionLunarCharge:
		reactionMultiplier = 3
	case attacks.AttackTagReactionLunarCrystallize:
		reactionMultiplier = 1.6
	}
	return (reactionMultiplier*(1+((6*em)/(2000+em))+reactBonus)*CalcReactionBaseDmg(lvl)*(1+atk.BaseDmgBonus) + atk.FlatDmg) * (1 + atk.Elevation)
}

func CalcReactionDmg(lvl int, src reactionBonusSrc, atk info.AttackInfo, em float64) (float64, info.Snapshot) {
	snap := info.Snapshot{
		CharLvl: lvl,
	}
	snap.Stats[attributes.EM] = em
	return (1 + ((16 * em) / (2000 + em)) + src.ReactBonus(atk)) * CalcReactionBaseDmg(lvl), snap
}

func CalcCatalyzeDmg(lvl int, src reactionBonusSrc, atk info.AttackInfo, em float64) float64 {
	return (1 + ((5 * em) / (1200 + em)) + src.ReactBonus(atk)) * CalcReactionBaseDmg(lvl)
}
