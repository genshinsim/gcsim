package aura

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	auraUptime []map[string]*calc.StreamStats
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		auraUptime: make([]map[string]*calc.StreamStats, len(cfg.Targets)),
	}

	for i := 0; i < len(cfg.Targets); i++ {
		out.auraUptime[i] = make(map[string]*calc.StreamStats)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i, c := range result.Enemies {
		for k, v := range c.ReactionUptime {
			if _, ok := b.auraUptime[i][k]; !ok {
				b.auraUptime[i][k] = &calc.StreamStats{}
			}
			b.auraUptime[i][k].Add(float64(v) / float64(result.Duration) * 100)
		}
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.TargetAuraUptime = make([]*model.SourceStats, len(b.auraUptime))
	for i, c := range b.auraUptime {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.TargetAuraUptime[i] = &model.SourceStats{
			Sources: source,
		}
	}
}
