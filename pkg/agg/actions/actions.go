package actions

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
	characterActions []map[string]*calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		characterActions: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.characterActions[i] = make(map[string]*calc.StreamStats)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i, c := range result.Characters {
		characterActions := make(map[string]float64)
		for _, ev := range c.ActionEvents {
			characterActions[ev.Action] += 1
		}
		for k, v := range characterActions {
			if _, ok := b.characterActions[i][k]; !ok {
				b.characterActions[i][k] = &calc.StreamStats{}
			}
			b.characterActions[i][k].Add(v)
		}
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.CharacterActions = make([]*model.SourceStats, len(b.characterActions))
	for i, c := range b.characterActions {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.CharacterActions[i] = &model.SourceStats{
			Sources: source,
		}
	}
}
