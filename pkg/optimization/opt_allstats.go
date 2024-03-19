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
//
// We also check if losing an ER sub and gaining a DMG sub is an overall gain. This will cover when the initial ER heuristic fails
// due to if .char.burst.ready { char burst; } lines not being modified as recommended.
func (stats *SubstatOptimizerDetails) optimizeERAndDMGSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	relevantSubstats := stats.charRelevantSubstats[idxChar]

	totalSubs := stats.getCharSubstatTotal(idxChar)
	if totalSubs != stats.totalLiquidSubstats {
		opDebug = append(opDebug, fmt.Sprint("Character has", totalSubs, "total liquid subs allocated but expected", stats.totalLiquidSubstats))
	}

	addedEr := false
	// Check if adding ER subs adds damage
	for stats.charMaxExtraERSubs[idxChar] > 0.0 && stats.charSubstatFinal[idxChar][attributes.ER] < stats.charSubstatLimits[idxChar][attributes.ER] {
		origIter := stats.simcfg.Settings.Iterations
		stats.simcfg.Settings.IgnoreBurstEnergy = true
		stats.simcfg.Settings.Iterations = 25
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, -1)
		stats.simcfg.Settings.IgnoreBurstEnergy = false
		stats.simcfg.Settings.Iterations = 200
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

		// If the overall damage gain is less than 0, we are done
		if erGainGradient[0]+lowestLoss <= 0 || lowestSub == attributes.NoStat {
			break
		}
		addedEr = true

		stats.charSubstatFinal[idxChar][lowestSub] -= 1
		stats.charProfilesCopy[idxChar].Stats[lowestSub] -= float64(1) * stats.substatValues[lowestSub] * stats.charSubstatRarityMod[idxChar]

		stats.charSubstatFinal[idxChar][attributes.ER] += 1
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] += float64(1) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
		stats.charMaxExtraERSubs[idxChar] -= 1
	}

	// Check if removing ER subs adds damage
	// We use less iterations and a higher threshold here because we prefer having more ER for better stddev.
	// This should help cover the case where users don't apply the suggestions to modify .char.burst.ready into
	// their intended rotation change. This still might get stuck in a local minima.
	// TODO: How to adjust for local minima?
	for !addedEr && stats.charSubstatFinal[idxChar][attributes.ER] > 0 {
		origIter := stats.simcfg.Settings.Iterations
		stats.simcfg.Settings.IgnoreBurstEnergy = true
		stats.simcfg.Settings.Iterations = 25
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, 1)
		stats.simcfg.Settings.IgnoreBurstEnergy = false
		stats.simcfg.Settings.Iterations = 100
		erGainGradient := stats.calculateSubstatGradientsForChar(idxChar, []attributes.Stat{attributes.ER}, -1)
		stats.simcfg.Settings.Iterations = origIter
		largestGain := -999999999999.0
		largestSub := attributes.NoStat
		for idxSubstat, gradient := range substatGradients {
			substat := relevantSubstats[idxSubstat]
			if stats.charSubstatFinal[idxChar][substat] < stats.charSubstatLimits[idxChar][substat] && gradient > largestGain {
				largestGain = gradient
				largestSub = substat
			}
		}

		// If the overall damage gain is less than 100, we are done
		if erGainGradient[0]+largestGain <= 100 || largestSub == attributes.NoStat {
			break
		}

		stats.charSubstatFinal[idxChar][largestSub] += 1
		stats.charProfilesCopy[idxChar].Stats[largestSub] += float64(1) * stats.substatValues[largestSub] * stats.charSubstatRarityMod[idxChar]

		stats.charSubstatFinal[idxChar][attributes.ER] -= 1
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] -= float64(1) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
		stats.charMaxExtraERSubs[idxChar] += 1
	}

	opDebug = append(opDebug, "Final "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	return opDebug
}
