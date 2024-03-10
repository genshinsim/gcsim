package endstats

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "ending-stats",
		New:  NewAgg,
	})
}

type buffer struct {
	endStats []endStats
}

type endStats struct {
	endingEnergy *calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		endStats: make([]endStats, len(cfg.Characters)),
	}

	for i := range out.endStats {
		out.endStats[i] = endStats{
			endingEnergy: &calc.StreamStats{},
		}
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i, v := range result.EndStats {
		b.endStats[i].endingEnergy.Add(v.EndingEnergy)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.EndStats = make([]*model.EndStats, len(b.endStats))
	for i, v := range b.endStats {
		result.EndStats[i] = &model.EndStats{
			EndingEnergy: agg.ToDescriptiveStats(v.endingEnergy),
		}
	}
}
