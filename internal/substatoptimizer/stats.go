package substatoptimizer

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type OptimStats struct {
	charRelevantSubstats map[keys.Char][]attributes.Stat
	artifactSets4Star    []keys.Set
	substatValues        []float64
	mainstatValues       []float64
}

func InitOptimStats() *OptimStats {
	c := OptimStats{}

	c.artifactSets4Star = []keys.Set{
		keys.TheExile,
		keys.Instructor,
	}

	c.substatValues = make([]float64, attributes.EndStatType)
	c.mainstatValues = make([]float64, attributes.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	c.substatValues[attributes.ATKP] = 0.0496
	c.substatValues[attributes.CR] = 0.0331
	c.substatValues[attributes.CD] = 0.0662
	c.substatValues[attributes.EM] = 19.82
	c.substatValues[attributes.ER] = 0.0551
	c.substatValues[attributes.HPP] = 0.0496
	c.substatValues[attributes.DEFP] = 0.062
	c.substatValues[attributes.ATK] = 16.54
	c.substatValues[attributes.DEF] = 19.68
	c.substatValues[attributes.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	c.mainstatValues[attributes.ATKP] = 0.466
	c.mainstatValues[attributes.CR] = 0.311
	c.mainstatValues[attributes.CD] = 0.622
	c.mainstatValues[attributes.EM] = 186.5
	c.mainstatValues[attributes.ER] = 0.518
	c.mainstatValues[attributes.HPP] = 0.466
	c.mainstatValues[attributes.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	c.charRelevantSubstats = map[keys.Char][]attributes.Stat{
		keys.Albedo:  {attributes.DEFP},
		keys.Hutao:   {attributes.HPP},
		keys.Kokomi:  {attributes.HPP},
		keys.Zhongli: {attributes.HPP},
		keys.Itto:    {attributes.DEFP},
		keys.Yunjin:  {attributes.DEFP},
		keys.Noelle:  {attributes.DEFP},
		keys.Gorou:   {attributes.DEFP},
	}

	return &c
}
