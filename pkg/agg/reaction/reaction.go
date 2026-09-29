package reaction

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	agg.Register(agg.Config{
		Name: "reaction",
		New:  NewAgg,
	})
}

type buffer struct {
	sourceReactions []map[string]*calc.StreamStats
	iters           uint
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		sourceReactions: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.sourceReactions[i] = make(map[string]*calc.StreamStats)
	}

	return &out, nil
}

func (b *buffer) Add(result *agg.Summary) {
	for i := range result.Characters {
		for k, v := range result.Characters[i].Reactions {
			if _, ok := b.sourceReactions[i][k]; !ok {
				b.sourceReactions[i][k] = &calc.StreamStats{}
			}
			b.sourceReactions[i][k].Add(v)
		}
	}
	b.iters++
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.SourceReactions = make([]*model.SourceStats, len(b.sourceReactions))
	for i, c := range b.sourceReactions {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			agg.PadStreamStatToCount(s, b.iters)
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.SourceReactions[i] = &model.SourceStats{
			Sources: source,
		}
	}
}
