package shield

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

const normalized string = "normalized"

var types = [...]string{
	normalized,
	attributes.Anemo.String(),
	attributes.Cryo.String(),
	attributes.Electro.String(),
	attributes.Geo.String(),
	attributes.Hydro.String(),
	attributes.Pyro.String(),
	attributes.Dendro.String(),
	attributes.Physical.String(),
}

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	shieldHP map[string]map[string]*stats.WeightedStreamStats
	uptime   map[string]*calc.StreamStats
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		shieldHP: make(map[string]map[string]*stats.WeightedStreamStats),
		uptime:   make(map[string]*calc.StreamStats),
	}

	out.shieldHP["effective"] = make(map[string]*stats.WeightedStreamStats)
	for _, t := range types {
		out.shieldHP["effective"][t] = &stats.WeightedStreamStats{}
	}
	out.uptime["effective"] = &calc.StreamStats{}
	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for _, shield := range result.ShieldResults.Shields {
		// create empty state if new shield
		if _, ok := b.shieldHP[shield.Name]; !ok {
			b.shieldHP[shield.Name] = make(map[string]*stats.WeightedStreamStats)
			for _, t := range types {
				b.shieldHP[shield.Name][t] = &stats.WeightedStreamStats{}
			}
			b.uptime[shield.Name] = &calc.StreamStats{}
		}

		var shieldUptime float64
		for _, interval := range shield.Intervals {
			weight := clamp(interval.End, result.Duration) - interval.Start
			shieldUptime += float64(weight)
			for k, hp := range interval.HP {
				b.shieldHP[shield.Name][k].Add(hp, weight)
			}
		}
		b.uptime[shield.Name].Add(shieldUptime / float64(result.Duration))
	}

	var effectiveUptime float64
	for k, intervals := range result.ShieldResults.EffectiveShield {
		for _, interval := range intervals {
			weight := clamp(interval.End, result.Duration) - interval.Start
			b.shieldHP["effective"][k].Add(interval.HP, weight)
			if k == normalized {
				effectiveUptime += float64(weight)
			}
		}
	}
	b.uptime["effective"].Add(effectiveUptime / float64(result.Duration))
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.Shields = make(map[string]*model.ShieldInfo)
	for k, s := range b.shieldHP {
		outHP := make(map[string]*model.DescriptiveStats)
		for t, hp := range s {
			mean, std := hp.Mean(), hp.StdDev()
			outHP[t] = &model.DescriptiveStats{
				Min:  &hp.Min,
				Max:  &hp.Max,
				Mean: &mean,
				SD:   &std,
			}
		}

		result.Shields[k] = &model.ShieldInfo{
			Hp:     outHP,
			Uptime: agg.ToDescriptiveStats(b.uptime[k]),
		}
	}
}

func clamp(x, max int) int {
	if x > max {
		return max
	}
	return x
}
