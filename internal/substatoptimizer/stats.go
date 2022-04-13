package substatoptimizer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"math"
)

type OptimStats struct {
	charRelevantSubstats  map[core.CharKey][]core.StatType
	artifactSets4Star     []string
	substatValues         []float64
	mainstatValues        []float64
	charSubstatFinal      [][]int
	charSubstatLimits     [][]int
	charSubstatRarityMod  []float64
	indivSubstatLiquidCap int
	simcfg                core.SimulationConfig
}

func InitOptimStats(simcfg core.SimulationConfig, indivLiquidCap int) *OptimStats {
	s := OptimStats{}

	// TODO: Will need to update this once artifact keys are introduced, and if more 4* artifact sets are implemented
	s.artifactSets4Star = []string{
		"exile",
		"instructor",
		"theexile",
	}

	s.substatValues = make([]float64, core.EndStatType)
	s.mainstatValues = make([]float64, core.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	s.substatValues[core.ATKP] = 0.0496
	s.substatValues[core.CR] = 0.0331
	s.substatValues[core.CD] = 0.0662
	s.substatValues[core.EM] = 19.82
	s.substatValues[core.ER] = 0.0551
	s.substatValues[core.HPP] = 0.0496
	s.substatValues[core.DEFP] = 0.062
	s.substatValues[core.ATK] = 16.54
	s.substatValues[core.DEF] = 19.68
	s.substatValues[core.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	s.mainstatValues[core.ATKP] = 0.466
	s.mainstatValues[core.CR] = 0.311
	s.mainstatValues[core.CD] = 0.622
	s.mainstatValues[core.EM] = 186.5
	s.mainstatValues[core.ER] = 0.518
	s.mainstatValues[core.HPP] = 0.466
	s.mainstatValues[core.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	s.charRelevantSubstats = map[core.CharKey][]core.StatType{
		core.Albedo:  {core.DEFP},
		core.Hutao:   {core.HPP},
		core.Kokomi:  {core.HPP},
		core.Zhongli: {core.HPP},
		core.Itto:    {core.DEFP},
		core.Yunjin:  {core.DEFP},
		core.Noelle:  {core.DEFP},
		core.Gorou:   {core.DEFP},
	}

	// Final output array that holds [character][substat_count]
	s.charSubstatFinal = make([][]int, len(simcfg.Characters.Profile))
	for i := range simcfg.Characters.Profile {
		s.charSubstatFinal[i] = make([]int, core.EndStatType)
	}

	s.indivSubstatLiquidCap = indivLiquidCap
	s.charSubstatLimits = make([][]int, len(simcfg.Characters.Profile))
	s.charSubstatRarityMod = make([]float64, len(simcfg.Characters.Profile))

	return &s
}

// Obtain substat count limits based on main stats and also determine 4* set status
// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
// Most people will have 1 5* artifact which messes things up
// TODO: Check whether taking like an average of the two stat values is good enough?
func (stats *OptimStats) setStatLimits() bool {
	profileIncludesFourStar := false

	for i, char := range stats.simcfg.Characters.Profile {
		stats.charSubstatLimits[i] = make([]int, core.EndStatType)
		for idxStat, stat := range stats.mainstatValues {
			if stat == 0 {
				continue
			}
			if char.Stats[idxStat] == 0 {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap
			} else {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap - (2 * int(math.Round(char.Stats[idxStat]/stats.mainstatValues[idxStat])))
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
