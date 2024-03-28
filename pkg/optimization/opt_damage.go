package optimization

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// TODO: Seems like this should be configurable
func (stats *SubstatOptimizerDetails) getNonErSubstatsToOptimizeForChar(char info.CharacterProfile) []attributes.Stat {
	// Get relevant substats, and add additional ones for special characters if needed
	relevantSubstats := []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD, attributes.EM}
	// RIP crystallize...
	if keys.CharKeyToEle[char.Base.Key] == attributes.Geo {
		relevantSubstats = []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD}
	}
	return relevantSubstats
}

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

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}

	// start from max liquid in all relevant substats
	for _, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat]) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
		stats.charSubstatFinal[idxChar][substat] = stats.charSubstatLimits[idxChar][substat]
	}

	totalSubs := stats.getCharSubstatTotal(idxChar)
	stats.optimizer.logger.Debug(char.Base.Key.Pretty())
	for totalSubs > stats.charTotalLiquidSubstats[idxChar] {
		amount := -1
		if totalSubs-stats.charTotalLiquidSubstats[idxChar] >= 8 {
			// reduce 5 at a time to quickly go from 10 liquid in an useless sub to 0 liquid
			amount = -5
		} else if totalSubs-stats.charTotalLiquidSubstats[idxChar] >= 4 {
			amount = -2
		}
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, amount)
		allocDebug := stats.allocateSomeSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats, amount)
		opDebug = append(opDebug, allocDebug...)
		totalSubs = stats.getCharSubstatTotal(idxChar)
		stats.optimizer.logger.Debug(totalSubs, " Liquid Substat Counts: "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	}
	opDebug = append(opDebug, "Liquid Substat Counts: "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))

	return opDebug
}
