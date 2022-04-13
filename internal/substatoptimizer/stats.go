package substatoptimizer

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type OptimStats struct {
	charRelevantSubstats map[keys.Char][]attributes.Stat
	artifactSets4Star    []keys.Set
	substatValues         []float64
	mainstatValues        []float64
	charSubstatFinal      [][]int
	charSubstatLimits     [][]int
	charSubstatRarityMod  []float64
	indivSubstatLiquidCap int
	fixedSubstatCount			int
	simcfg                *ast.ActionList
}

func InitOptimStats(simcfg *ast.ActionList, indivLiquidCap int, fixedSubstatCount int) *OptimStats {
	s := OptimStats{}

	s.artifactSets4Star = []keys.Set{
		keys.TheExile,
		keys.Instructor,
	}

	s.substatValues = make([]float64, attributes.EndStatType)
	s.mainstatValues = make([]float64, attributes.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	s.substatValues[attributes.ATKP] = 0.0496
	s.substatValues[attributes.CR] = 0.0331
	s.substatValues[attributes.CD] = 0.0662
	s.substatValues[attributes.EM] = 19.82
	s.substatValues[attributes.ER] = 0.0551
	s.substatValues[attributes.HPP] = 0.0496
	s.substatValues[attributes.DEFP] = 0.062
	s.substatValues[attributes.ATK] = 16.54
	s.substatValues[attributes.DEF] = 19.68
	s.substatValues[attributes.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	s.mainstatValues[attributes.ATKP] = 0.466
	s.mainstatValues[attributes.CR] = 0.311
	s.mainstatValues[attributes.CD] = 0.622
	s.mainstatValues[attributes.EM] = 186.5
	s.mainstatValues[attributes.ER] = 0.518
	s.mainstatValues[attributes.HPP] = 0.466
	s.mainstatValues[attributes.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	s.charRelevantSubstats = map[keys.Char][]attributes.Stat{
		keys.Albedo:  {attributes.DEFP},
		keys.Hutao:   {attributes.HPP},
		keys.Kokomi:  {attributes.HPP},
		keys.Zhongli: {attributes.HPP},
		keys.Itto:    {attributes.DEFP},
		keys.Yunjin:  {attributes.DEFP},
		keys.Noelle:  {attributes.DEFP},
		keys.Gorou:   {attributes.DEFP},
	}

	// Final output array that holds [character][substat_count]
	s.charSubstatFinal = make([][]int, len(simcfg.Characters))
	for i := range simcfg.Characters {
		s.charSubstatFinal[i] = make([]int, attributes.EndStatType)
	}

	s.simcfg = simcfg
	s.indivSubstatLiquidCap = indivLiquidCap
	s.fixedSubstatCount = fixedSubstatCount
	s.charSubstatLimits = make([][]int, len(simcfg.Characters))
	s.charSubstatRarityMod = make([]float64, len(simcfg.Characters))

	return &s
}

// Obtain substat count limits based on main stats and also determine 4* set status
// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
// Most people will have 1 5* artifact which messes things up
// TODO: Check whether taking like an average of the two stat values is good enough?
func (stats *OptimStats) setStatLimits() bool {
	profileIncludesFourStar := false

	for i, char := range stats.simcfg.Characters {
		stats.charSubstatLimits[i] = make([]int, attributes.EndStatType)
		for idxStat, stat := range stats.mainstatValues {
			if stat == 0 {
				continue
			}
			if char.Stats[idxStat] == 0 {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap
			} else {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap - (stats.fixedSubstatCount * int(math.Round(char.Stats[idxStat]/stats.mainstatValues[idxStat])))
			}
		}

		// Display warning message for 4* sets
		stats.charSubstatRarityMod[i] = 1
		for set := range char.Sets {
			for _, fourStar := range stats.artifactSets4Star {
				if set == fourStar {
					profileIncludesFourStar = true
					stats.charSubstatRarityMod[i] = 0.8
				}
			}
		}
	}

	return profileIncludesFourStar
}
