package overview

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "overview",
		New:  NewAgg,
	})
}

type buffer struct {
	duration    *calc.Sample
	dps         *calc.Sample
	rps         *calc.Sample
	eps         *calc.Sample
	hps         *calc.Sample
	shp         *calc.Sample
	totalDamage calc.StreamStats
}

func newSample(itr int) *calc.Sample {
	return &calc.Sample{
		Xs:     make([]float64, 0, itr),
		Sorted: false,
	}
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		duration:    newSample(cfg.Settings.Iterations),
		dps:         newSample(cfg.Settings.Iterations),
		rps:         newSample(cfg.Settings.Iterations),
		hps:         newSample(cfg.Settings.Iterations),
		eps:         newSample(cfg.Settings.Iterations),
		shp:         newSample(cfg.Settings.Iterations),
		totalDamage: calc.StreamStats{},
	}
	return &out, nil
}

// TODO: push looping/summation to StatsCollector for peformance boost
func (b *buffer) Add(result stats.Result) {
	b.totalDamage.Add(result.TotalDamage)

	b.duration.Xs = append(b.duration.Xs, float64(result.Duration)/60)
	b.dps.Xs = append(b.dps.Xs, result.DPS)
	b.rps.Xs = append(b.rps.Xs, 0)
	b.hps.Xs = append(b.hps.Xs, 0)
	b.eps.Xs = append(b.eps.Xs, 0)
	b.shp.Xs = append(b.shp.Xs, 0)
	iX := len(b.shp.Xs) - 1

	for _, interval := range result.ShieldResults.EffectiveShield["normalized"] {
		end := interval.End
		if end > result.Duration {
			end = result.Duration
		}
		b.shp.Xs[iX] += interval.HP * float64(end-interval.Start)
	}
	b.shp.Xs[iX] /= float64(result.Duration)

	for i := range result.Characters {
		b.rps.Xs[iX] += float64(len(result.Characters[i].ReactionEvents))
		for _, h := range result.Characters[i].HealEvents {
			b.hps.Xs[iX] += h.Heal
		}
		for _, e := range result.Characters[i].EnergyEvents {
			b.eps.Xs[iX] += e.Gained + e.Wasted
		}
	}

	b.rps.Xs[iX] /= b.duration.Xs[iX]
	b.hps.Xs[iX] /= b.duration.Xs[iX]
	b.eps.Xs[iX] /= b.duration.Xs[iX]
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.Duration = agg.ToOverviewStats(b.duration)
	result.DPS = agg.ToOverviewStats(b.dps)
	result.RPS = agg.ToOverviewStats(b.rps)
	result.EPS = agg.ToOverviewStats(b.eps)
	result.HPS = agg.ToOverviewStats(b.hps)
	result.SHP = agg.ToOverviewStats(b.shp)
	result.TotalDamage = agg.ToDescriptiveStats(&b.totalDamage)
}
