package energycalc

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
	ercalc     bool
	erNeeded   []*calc.Sample
	weightedER []*calc.Sample
}

func newSample(itr int) *calc.Sample {
	return &calc.Sample{
		Xs:     make([]float64, 0, itr),
		Sorted: false,
	}
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	if !cfg.Settings.ErCalc {
		out := buffer{
			ercalc: false,
		}
		return &out, nil
	}
	out := buffer{
		ercalc:     true,
		erNeeded:   make([]*calc.Sample, len(cfg.Characters)),
		weightedER: make([]*calc.Sample, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.erNeeded[i] = newSample(cfg.Settings.Iterations)
		out.weightedER[i] = newSample(cfg.Settings.Iterations)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	if !b.ercalc {
		return
	}
	for i := range result.Characters {
		b.erNeeded[i].Xs = append(b.erNeeded[i].Xs, result.Characters[i].ErNeeded)
		b.weightedER[i].Xs = append(b.weightedER[i].Xs, result.Characters[i].WeightedER)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	if !b.ercalc {
		return
	}
	result.ErNeeded = make([]*model.OverviewStats, len(b.erNeeded))
	result.WeightedEr = make([]*model.OverviewStats, len(b.weightedER))
	for i, c := range b.erNeeded {
		result.ErNeeded[i] = agg.ToOverviewStats(c)
	}
	for i, c := range b.weightedER {
		result.WeightedEr[i] = agg.ToOverviewStats(c)
	}
}
