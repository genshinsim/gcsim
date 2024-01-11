package optimization

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// Calculate per-character per-substat "gradients" at initial state using finite differences
func (stats *SubstatOptimizerDetails) optimizeERAndDMGSubstats() []string {
	var (
		opDebug   []string
		charDebug []string
	)

	stats.simcfg.Characters = stats.charProfilesCopy

	for idxChar := range stats.charProfilesCopy {
		charDebug = stats.optimizeERAndDMGSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar])
		opDebug = append(opDebug, charDebug...)
	}
	return opDebug
}

// This function assumes that there are now all subs allocated. For every sub of ER we gain, we will lose one sub of damage
// We compare the damage loss of 1 DMG sub against the damage gain of 1 ER sub.
// We deallocate that DMG sub and allocate 1 ER sub if it would be an overall gain
// Repeat until we cannot allocate ER subs or the DMG loss would be greater than the gain
// The ER dmg gain is prone to noise, so we need to do more iterations
func (stats *SubstatOptimizerDetails) optimizeERAndDMGSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}
	totalSubs := stats.getCharSubstatTotal(idxChar)
	if totalSubs != stats.totalLiquidSubstats {
		opDebug = append(opDebug, fmt.Sprint("Character has", totalSubs, "total liquid subs allocated but expected", stats.totalLiquidSubstats))
	}
	// fmt.Println(char.Base.Key.Pretty(), "has", totalSubs, "total liquid substats")
	for stats.charMaxExtraERSubs[idxChar] > 0.0 && stats.charSubstatFinal[idxChar][attributes.ER] < stats.charSubstatLimits[idxChar][attributes.ER] {
		origIter := stats.simcfg.Settings.Iterations
		stats.simcfg.Settings.IgnoreBurstEnergy = true
		stats.simcfg.Settings.Iterations = 25
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, -1)
		stats.simcfg.Settings.IgnoreBurstEnergy = false
		stats.simcfg.Settings.Iterations = 250
		erGainGradient := stats.calculateSubstatGradientsForChar(idxChar, []attributes.Stat{attributes.ER}, 1)
		stats.simcfg.Settings.Iterations = origIter
		lowestLoss := -999999999999.0
		lowestSub := attributes.NoStat
		for idxSubstat, gradient := range substatGradients {
			substat := relevantSubstats[idxSubstat]
			if stats.charSubstatFinal[idxChar][substat] > 0 && gradient > lowestLoss {
				lowestLoss = gradient
				lowestSub = substat
			}
		}
		if erGainGradient[0]+lowestLoss <= 0 || lowestSub == attributes.NoStat {
			break
		}

		stats.charSubstatFinal[idxChar][lowestSub] -= 1
		stats.charProfilesCopy[idxChar].Stats[lowestSub] -= float64(1) * stats.substatValues[lowestSub] * stats.charSubstatRarityMod[idxChar]

		stats.charSubstatFinal[idxChar][attributes.ER] += 1
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] += float64(1) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
		stats.charMaxExtraERSubs[idxChar] -= 1
	}
	opDebug = append(opDebug, "Final Liquid Substat Counts: "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	return opDebug
}
