package optimization

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// Calculate per-character per-substat "gradients" at initial state using finite differences
// We use ignore_burst_energy mode to remove noise from energy, and the custom damage collector
// is used to remove noise from random crit. This allows us to run a very small 25 iterations per gradient calculation
//
// TODO: Add setting which allows the user to increase the number of iterations (for cases
// with inherent randomness like Widsith or random delays)
//
// TODO: Automatically increase iteration count when stddev is very high?
func (stats *SubstatOptimizerDetails) optimizeNonERSubstats() []string {
	var (
		opDebug   []string
		charDebug []string
	)
	origIter := stats.simcfg.Settings.Iterations
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	stats.simcfg.Settings.Iterations = 25
	stats.simcfg.Characters = stats.charProfilesCopy

	for idxChar := range stats.charProfilesCopy {
		charDebug = stats.optimizeNonErSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar])
		opDebug = append(opDebug, charDebug...)
	}
	stats.simcfg.Settings.IgnoreBurstEnergy = false
	stats.simcfg.Settings.Iterations = origIter
	return opDebug
}

// This calculation starts all the relevant substats at maximum allocated liquid.
// This reduces the chances of hitting a local maximum of stacking atk%/hp%/def%
// It uses a gradient to determine which substat would cause the least damage loss
// when removing. This continues until we are within the total liquid limits
// Initially substats are removed 5 or 2 at a time to speed up the computations
//
// TODO: Allow the user to specify the removal rate?
//
// TODO: Multistart gradient descent/ascent from 0 allocated liquid and compare?
func (stats *SubstatOptimizerDetails) optimizeNonErSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	// Reset favonius char crit rate
	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[attributes.CR] -= FavCritRateBias * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
	}

	var relevantSubstats []attributes.Stat
	relevantSubstats = append(relevantSubstats, stats.charRelevantSubstats[idxChar]...)

	// start from max liquid in all relevant substats
	for _, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat]) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
		stats.charSubstatFinal[idxChar][substat] = stats.charSubstatLimits[idxChar][substat]
	}

	totalSubs := stats.getCharSubstatTotal(idxChar)
	stats.optimizer.logger.Debug(char.Base.Key.Pretty())
	stats.optimizer.logger.Debug(PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	for totalSubs > stats.totalLiquidSubstats {
		amount := -1
		switch {
		case totalSubs-stats.totalLiquidSubstats >= 15:
			amount = -20 // will get clamped to either 10/8 depending on the substat limit
		case totalSubs-stats.totalLiquidSubstats >= 8:
			amount = -5
		case totalSubs-stats.totalLiquidSubstats >= 4:
			amount = -2
		}
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, amount)

		// loops multiple gradients while totalSubs-stats.totalLiquidSubstats >= 25
		// this should be most correct because the first 5 to 6 substats have 0 effect on dps
		for ok := true; ok; ok = totalSubs-stats.totalLiquidSubstats >= 25 {
			allocDebug := stats.allocateSomeSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats, amount)
			totalSubs = stats.getCharSubstatTotal(idxChar)
			opDebug = append(opDebug, allocDebug...)

			// filter out substats that are at minimum
			newRelevantSubstats := []attributes.Stat{}
			newSubstatGrad := []float64{}
			removedGrad := -100000000.0
			for idxSub, substat := range relevantSubstats {
				if stats.charSubstatFinal[idxChar][substat] > 0 {
					newRelevantSubstats = append(newRelevantSubstats, substat)
					newSubstatGrad = append(newSubstatGrad, substatGradients[idxSub])
				} else {
					removedGrad = max(removedGrad, substatGradients[idxSub])
				}
			}
			// only update the charRelevantSubstats when the gradient of the removed substats is very small
			// this is used later in the opt_allstats
			if stats.getCharSubstatTotal(idxChar)-stats.totalLiquidSubstats >= 15 ||
				removedGrad >= -100 {
				stats.charRelevantSubstats[idxChar] = nil
				stats.charRelevantSubstats[idxChar] = append(stats.charRelevantSubstats[idxChar], newRelevantSubstats...)
			}
			relevantSubstats = newRelevantSubstats
			substatGradients = newSubstatGrad
		}
		stats.optimizer.logger.Debug(PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	}
	opDebug = append(opDebug, PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	stats.optimizer.logger.Debug(char.Base.Key, " has relevant substats:", stats.charRelevantSubstats[idxChar])
	return opDebug
}
