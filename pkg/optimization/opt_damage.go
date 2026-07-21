package optimization

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (stats *SubstatOptimizerDetails) optimizeNonERSubstats() []string {
	var opDebug []string
	origIter := stats.simcfg.Settings.Iterations
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	stats.simcfg.Settings.Iterations = 25
	stats.simcfg.Characters = stats.charProfilesCopy

	for idxChar := len(stats.charProfilesCopy) - 1; idxChar >= 0; idxChar-- {
		charDebug := stats.optimizeNonErSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar])
		opDebug = append(opDebug, charDebug...)
	}

	stats.simcfg.Settings.IgnoreBurstEnergy = false
	stats.simcfg.Settings.Iterations = origIter
	return opDebug
}

func (stats *SubstatOptimizerDetails) optimizeNonErSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[attributes.CR] -= FavCritRateBias * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
	}

	relevant := stats.charRelevantSubstats[idxChar]

	for _, s := range relevant {
		stats.charProfilesCopy[idxChar].Stats[s] -= float64(stats.charSubstatFinal[idxChar][s]) *
			stats.substatValues[s] * stats.charSubstatRarityMod[idxChar]
		stats.charSubstatFinal[idxChar][s] = 0
	}

	target := stats.charTotalLiquidSubstats[idxChar]

	origIter := stats.simcfg.Settings.Iterations
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	stats.simcfg.Settings.Iterations = 25

	gradients := stats.calculateSubstatGradientsForChar(idxChar, relevant, 1)

	maxGrad := 0.0
	for _, g := range gradients {
		if g > maxGrad {
			maxGrad = g
		}
	}
	if maxGrad > 0 {
		filtered := make([]attributes.Stat, 0, len(relevant))
		filteredGrads := make([]float64, 0, len(gradients))
		for i, g := range gradients {
			if g >= maxGrad*0.01 {
				filtered = append(filtered, relevant[i])
				filteredGrads = append(filteredGrads, g)
			}
		}
		relevant = filtered
		gradients = filteredGrads
	}

	totalSubs := 0
	for totalSubs < target {
		bestIdx := 0
		bestGain := -1e99
		for i, g := range gradients {
			if g > bestGain && stats.charSubstatFinal[idxChar][relevant[i]] < stats.charSubstatLimits[idxChar][relevant[i]] {
				bestGain = g
				bestIdx = i
			}
		}
		if bestGain <= 0 {
			break
		}
		best := relevant[bestIdx]

		stats.charSubstatFinal[idxChar][best]++
		stats.charProfilesCopy[idxChar].Stats[best] += stats.substatValues[best] * stats.charSubstatRarityMod[idxChar]
		totalSubs++

		gradients[bestIdx] = stats.calculateSingleSubstatGradient(idxChar, best, 1)
	}

	stats.simcfg.Settings.Iterations = origIter

	opDebug = append(opDebug, PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	stats.optimizer.logger.Debug(char.Base.Key, " has relevant substats:", stats.charRelevantSubstats[idxChar])
	return opDebug
}
