package energy

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	sourceEnergy []map[string]*calc.StreamStats
	erNeeded     []*calc.Sample
	weightedER   []*calc.Sample
}

func newSample(itr int) *calc.Sample {
	return &calc.Sample{
		Xs:     make([]float64, 0, itr),
		Sorted: false,
	}
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		sourceEnergy: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		erNeeded:     make([]*calc.Sample, len(cfg.Characters)),
		weightedER:   make([]*calc.Sample, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.sourceEnergy[i] = make(map[string]*calc.StreamStats)
		out.erNeeded[i] = newSample(cfg.Settings.Iterations)
		out.weightedER[i] = newSample(cfg.Settings.Iterations)
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
		b.erNeeded[i].Xs = append(b.erNeeded[i].Xs, result.Characters[i].ERneeded)
		b.weightedER[i].Xs = append(b.weightedER[i].Xs, result.Characters[i].WeightedER)
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

	result.ErNeeded = make([]*model.OverviewStats, len(b.erNeeded))
	result.WeightedEr = make([]*model.OverviewStats, len(b.erNeeded))
	for i, c := range b.erNeeded {
		result.ErNeeded[i] = agg.ToOverviewStats(c)
	}
	for i, c := range b.weightedER {
		result.WeightedEr[i] = agg.ToOverviewStats(c)
	}
}
