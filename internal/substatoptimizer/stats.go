package substatoptimizer

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type OptimStats struct {
	charRelevantSubstats map[core.CharKey][]core.StatType
	artifactSets4Star    []string
	substatValues        []float64
	mainstatValues       []float64
}

func InitOptimStats() *OptimStats {
	c := OptimStats{}

	// TODO: Will need to update this once artifact keys are introduced, and if more 4* artifact sets are implemented
	c.artifactSets4Star = []string{
		"exile",
		"instructor",
		"theexile",
	}

	c.substatValues = make([]float64, core.EndStatType)
	c.mainstatValues = make([]float64, core.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	c.substatValues[core.ATKP] = 0.0496
	c.substatValues[core.CR] = 0.0331
	c.substatValues[core.CD] = 0.0662
	c.substatValues[core.EM] = 19.82
	c.substatValues[core.ER] = 0.0551
	c.substatValues[core.HPP] = 0.0496
	c.substatValues[core.DEFP] = 0.062
	c.substatValues[core.ATK] = 16.54
	c.substatValues[core.DEF] = 19.68
	c.substatValues[core.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	c.mainstatValues[core.ATKP] = 0.466
	c.mainstatValues[core.CR] = 0.311
	c.mainstatValues[core.CD] = 0.622
	c.mainstatValues[core.EM] = 186.5
	c.mainstatValues[core.ER] = 0.518
	c.mainstatValues[core.HPP] = 0.466
	c.mainstatValues[core.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	c.charRelevantSubstats = map[core.CharKey][]core.StatType{
		core.Albedo:  {core.DEFP},
		core.Hutao:   {core.HPP},
		core.Kokomi:  {core.HPP},
		core.Zhongli: {core.HPP},
		core.Itto:    {core.DEFP},
		core.Yunjin:  {core.DEFP},
		core.Noelle:  {core.DEFP},
		core.Gorou:   {core.DEFP},
	}

	return &c
}
