package optimization

import (
	"context"
	"sort"
	"time"

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
	artifactSets4Star      []keys.Set
	substatValues          []float64
	mainstatValues         []float64
	charSubstatFinal       [][]int
	charSubstatLimits      [][]int
	charSubstatRarityMod   []float64
	charProfilesInitial    []info.CharacterProfile
	charWithFavonius       []bool
	charProfilesERBaseline []info.CharacterProfile
	charProfilesCopy       []info.CharacterProfile
	charMaxExtraERSubs     []float64
	charRelevantSubstats   [][]attributes.Stat
	simcfg                 *info.ActionList
	gcsl                   ast.Node
	simopt                 simulator.Options
	cfg                    string
	fixedSubstatCount      int
	indivSubstatLiquidCap  int
	totalLiquidSubstats    int
	optimizer              *SubstatOptimizer
}

func (stats *SubstatOptimizerDetails) allocateSomeSubstatGradientsForChar(
	idxChar int,
	_ info.CharacterProfile,
	substatGradient []float64,
	relevantSubstats []attributes.Stat,
	amount int,
) []string {
	var opDebug []string
	sorted := newSlice(substatGradient...)
	sort.Sort(sort.Reverse(sorted))

	for _, idxSubstat := range sorted.idx {
		substat := relevantSubstats[idxSubstat]

		if amount > 0 {
			if stats.charSubstatFinal[idxChar][substat] < stats.charSubstatLimits[idxChar][substat] {
				stats.charSubstatFinal[idxChar][substat] += amount
				stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
				return opDebug
			}
		}

		if stats.charSubstatFinal[idxChar][substat] > 0 {
			amount = clamp[int](-stats.charSubstatFinal[idxChar][substat], amount, amount)
			stats.charSubstatFinal[idxChar][substat] += amount
			stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
			return opDebug
		}
	}

	// TODO: No relevant substat can be allocated/deallocated, alloc/dealloc some random other substat??
	opDebug = append(opDebug, "Couldn't alloc/dealloc anything?????")
	return opDebug
}

func (stats *SubstatOptimizerDetails) calculateSubstatGradientsForChar(
	idxChar int,
	relevantSubstats []attributes.Stat,
	amount int,
) []float64 {
	stats.simcfg.Characters = stats.charProfilesCopy

	seed := time.Now().UnixNano()
	init := optstats.NewDamageAggBuffer(stats.simcfg)
	_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, init.Add)
	if err != nil {
		stats.optimizer.logger.Fatal(err.Error())
	}
	init.Flush()
	// TODO: Test if median or mean gives better results
	initialMean := mean(init.ExpectedDps)
	substatGradients := make([]float64, len(relevantSubstats))
	// Build "gradient" by substat
	for idxSubstat, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

		stats.simcfg.Characters = stats.charProfilesCopy

		a := optstats.NewDamageAggBuffer(stats.simcfg)
		_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerDmgStat, a.Add)
		if err != nil {
			stats.optimizer.logger.Fatal(err.Error())
		}
		a.Flush()

		substatGradients[idxSubstat] = mean(a.ExpectedDps) - initialMean
		// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
		// might give them "too much" cr (= max out liquid cr subs or overcap crit beyond 100%) but that's probably not a big deal
		if stats.simcfg.Settings.IgnoreBurstEnergy && stats.charWithFavonius[idxChar] && substat == attributes.CR {
			substatGradients[idxSubstat] += 1000 * float64(amount)
		}
		stats.charProfilesCopy[idxChar].Stats[substat] -= float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	}
	return substatGradients
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
			if attributes.Stat(idxStat) == attributes.ER {
				stats.charProfilesInitial[i].Stats[idxStat] += float64(fixedSubstatCount) * stat
			} else {
				stats.charProfilesInitial[i].Stats[idxStat] += float64(fixedSubstatCount) * stat * stats.charSubstatRarityMod[i]
			}
		}
	}
}
