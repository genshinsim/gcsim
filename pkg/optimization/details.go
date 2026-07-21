package optimization

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/optimization/optstats"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}
type SubstatOptimizerDetails struct {
	artifactSets4Star       []keys.Set
	substatValues           []float64
	mainstatValues          []float64
	mainstatTol             float64
	fourstarMod             float64
	charSubstatFinal        [][]int
	charSubstatLimits       [][]int
	charTotalLiquidSubstats []int
	charSubstatRarityMod    []float64
	charProfilesInitial     []info.CharacterProfile
	charWithFavonius        []bool
	charProfilesERBaseline  []info.CharacterProfile
	charRelevantSubstats    [][]attributes.Stat
	charProfilesCopy        []info.CharacterProfile
	charMaxExtraERSubs      []float64
	file                    *ast.File
	simcfg                  *info.ActionList
	gcsl                    ast.Node
	simopt                  simulator.Options
	cfg                     string
	fixedSubstatCount       int
	indivSubstatLiquidCap   int
	totalLiquidSubstats     int
	optimizer               *SubstatOptimizer
	gradientSeed            int64
}

// #1: deterministic seed — uses stats.gradientSeed which increments each call,
// but is constant within a batch so all N stats share the same noise.
func (stats *SubstatOptimizerDetails) calculateSubstatGradientsForChar(
	idxChar int,
	relevantSubstats []attributes.Stat,
	amount int,
) []float64 {
	stats.simcfg.Characters = stats.charProfilesCopy

	seed := stats.gradientSeed
	stats.gradientSeed++
	init := optstats.NewDamageAggBuffer(stats.simcfg)
	_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.file, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, init.Add)
	if err != nil {
		stats.optimizer.logger.Fatal(err.Error())
	}
	init.Flush()
	initialMean := mean(init.ExpectedDps)
	substatGradients := make([]float64, len(relevantSubstats))
	for idxSubstat, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

		stats.simcfg.Characters = stats.charProfilesCopy

		a := optstats.NewDamageAggBuffer(stats.simcfg)
		_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.file, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, a.Add)
		if err != nil {
			stats.optimizer.logger.Fatal(err.Error())
		}
		a.Flush()

		substatGradients[idxSubstat] = mean(a.ExpectedDps) - initialMean
		if stats.simcfg.Settings.IgnoreBurstEnergy && stats.charWithFavonius[idxChar] && substat == attributes.CR {
			substatGradients[idxSubstat] += 1000 * float64(amount)
		}
		stats.charProfilesCopy[idxChar].Stats[substat] -= float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	}
	return substatGradients
}

// #3: single-stat gradient for lazy updates — only recompute the stat that changed.
func (stats *SubstatOptimizerDetails) calculateSingleSubstatGradient(
	idxChar int,
	substat attributes.Stat,
	amount int,
) float64 {
	seed := stats.gradientSeed
	stats.gradientSeed++

	stats.simcfg.Characters = stats.charProfilesCopy
	init := optstats.NewDamageAggBuffer(stats.simcfg)
	_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.file, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, init.Add)
	if err != nil {
		stats.optimizer.logger.Fatal(err.Error())
	}
	init.Flush()
	initialMean := mean(init.ExpectedDps)

	stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	stats.simcfg.Characters = stats.charProfilesCopy

	a := optstats.NewDamageAggBuffer(stats.simcfg)
	_, err = optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.file, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, a.Add)
	if err != nil {
		stats.optimizer.logger.Fatal(err.Error())
	}
	a.Flush()

	gradient := mean(a.ExpectedDps) - initialMean
	if stats.simcfg.Settings.IgnoreBurstEnergy && stats.charWithFavonius[idxChar] && substat == attributes.CR {
		gradient += 1000 * float64(amount)
	}
	stats.charProfilesCopy[idxChar].Stats[substat] -= float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	return gradient
}

func (stats *SubstatOptimizerDetails) setInitialSubstats(fixedSubstatCount int) {
	stats.cloneStatsWithFixedAllocations(fixedSubstatCount)
	stats.calculateERBaseline()
}

// Copy to save initial character state with fixed allocations (2 of each substat)
func (stats *SubstatOptimizerDetails) cloneStatsWithFixedAllocations(fixedSubstatCount int) {
	for i := range stats.simcfg.Characters {
		stats.charProfilesInitial[i] = stats.simcfg.Characters[i].Clone()
		for idxStat, stat := range stats.substatValues {
			if stat == 0 {
				continue
			}
			stats.charProfilesInitial[i].Stats[idxStat] += float64(fixedSubstatCount) * stat * stats.charSubstatRarityMod[i]
		}
	}
}
