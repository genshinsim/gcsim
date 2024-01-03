package reaction

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "reaction",
		New:  NewAgg,
	})
}

type buffer struct {
	sourceReactions []map[string]*calc.StreamStats
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

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		sourceReactions := make(map[string]float64)
		for _, ev := range result.Characters[i].ReactionEvents {
			sourceReactions[ev.Reaction] += 1
		}
		for k, v := range sourceReactions {
			if _, ok := b.sourceReactions[i][k]; !ok {
				b.sourceReactions[i][k] = &calc.StreamStats{}
			}
			b.sourceReactions[i][k].Add(v)
		}
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.SourceReactions = make([]*model.SourceStats, len(b.sourceReactions))
	for i, c := range b.sourceReactions {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.SourceReactions[i] = &model.SourceStats{
			Sources: source,
		}
	}
}
