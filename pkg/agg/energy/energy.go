package energy

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "energy",
		New:  NewAgg,
	})
}

type buffer struct {
	sourceEnergy []map[string]*calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		sourceEnergy: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.sourceEnergy[i] = make(map[string]*calc.StreamStats)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		sourceEnergy := make(map[string]float64)
		for _, ev := range result.Characters[i].EnergyEvents {
			sourceEnergy[ev.Source] += ev.Gained + ev.Wasted
		}
		for k, v := range sourceEnergy {
			if _, ok := b.sourceEnergy[i][k]; !ok {
				b.sourceEnergy[i][k] = &calc.StreamStats{}
			}
			b.sourceEnergy[i][k].Add(v)
		}
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.TotalSourceEnergy = make([]*model.SourceStats, len(b.sourceEnergy))
	for i, c := range b.sourceEnergy {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.TotalSourceEnergy[i] = &model.SourceStats{
			Sources: source,
		}
	}
}
