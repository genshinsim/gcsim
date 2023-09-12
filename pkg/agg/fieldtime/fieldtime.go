package fieldtime

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
	fieldTimes []*calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		fieldTimes: make([]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.fieldTimes[i] = &calc.StreamStats{}
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		b.fieldTimes[i].Add(float64(result.Characters[i].ActiveTime) / 60)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.FieldTime = make([]*model.DescriptiveStats, len(b.fieldTimes))
	for i, c := range b.fieldTimes {
		result.FieldTime[i] = agg.ToDescriptiveStats(c)
	}
}
