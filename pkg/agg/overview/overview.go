package overview

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	duration    calc.Sample
	dps         calc.Sample
	rps         calc.Sample
	eps         calc.Sample
	hps         calc.Sample
	sps         calc.Sample
	totalDamage calc.StreamStats
}

func newSample(itr int) calc.Sample {
	return calc.Sample{
		Xs:     make([]float64, 0, itr),
		Sorted: false,
	}
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		duration:    newSample(cfg.Settings.Iterations),
		dps:         newSample(cfg.Settings.Iterations),
		rps:         newSample(cfg.Settings.Iterations),
		hps:         newSample(cfg.Settings.Iterations),
		eps:         newSample(cfg.Settings.Iterations),
		sps:         newSample(cfg.Settings.Iterations),
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
	b.sps.Xs = append(b.sps.Xs, 0)
	i := len(b.sps.Xs) - 1

	for _, s := range result.Shields {
		b.sps.Xs[i] += s.Absorption * float64(s.End-s.Start) / 60
	}

	for _, c := range result.Characters {
		b.rps.Xs[i] += float64(len(c.ReactionEvents))
		for _, h := range c.HealEvents {
			b.hps.Xs[i] += h.Heal
		}
		for _, e := range c.EnergyEvents {
			b.eps.Xs[i] += e.Gained + e.Wasted
		}
	}

	b.rps.Xs[i] /= b.duration.Xs[i]
	b.hps.Xs[i] /= b.duration.Xs[i]
	b.eps.Xs[i] /= b.duration.Xs[i]
	b.sps.Xs[i] /= b.duration.Xs[i]
}

func (b *buffer) Flush(result *agg.Result) {
	result.Duration = convert(b.duration)
	result.DPS = convert(b.dps)
	result.RPS = convert(b.rps)
	result.EPS = convert(b.eps)
	result.HPS = convert(b.hps)
	result.SPS = convert(b.sps)

	result.TotalDamage = agg.FloatStat{
		Min:  b.totalDamage.Min,
		Max:  b.totalDamage.Max,
		Mean: b.totalDamage.Mean(),
		SD:   b.totalDamage.StdDev(),
	}
}

func convert(input calc.Sample) agg.FloatStat {
	input.Sorted = false
	input.Sort()

	out := agg.FloatStat{
		Mean: input.Mean(),
		SD:   input.StdDev(),
		Q1:   input.Quantile(0.25),
		Q2:   input.Quantile(0.5),
		Q3:   input.Quantile(0.75),
	}
	out.Min, out.Max = input.Bounds()

	return out
}
