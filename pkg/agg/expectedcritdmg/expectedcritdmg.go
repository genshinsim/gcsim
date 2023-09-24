package expectedcritdmg

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
	expectedCritDmgMode  bool
	totalExpectedDPS     *calc.StreamStats
	characterExpectedDPS []*calc.StreamStats // i = char
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	if !cfg.Settings.ExpectedCritDmg {
		out := buffer{
			expectedCritDmgMode: false,
		}
		return &out, nil
	}

	out := buffer{
		expectedCritDmgMode:  true,
		totalExpectedDPS:     &calc.StreamStats{},
		characterExpectedDPS: make([]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.characterExpectedDPS[i] = &calc.StreamStats{}
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	if !b.expectedCritDmgMode {
		return
	}
	totalDPS := 0.0
	for i := range result.Characters {
		charDPS := 0.0
		if result.Duration > 0 {
			charDPS = result.Characters[i].ExpectedCritDmg / (float64(result.Duration) / 60.0)
		}
		b.characterExpectedDPS[i].Add(charDPS)
		totalDPS += charDPS
	}
	b.totalExpectedDPS.Add(totalDPS)
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	if !b.expectedCritDmgMode {
		return
	}
	result.CharacterExpectedDps = make([]*model.DescriptiveStats, len(b.characterExpectedDPS))
	for i := range b.characterExpectedDPS {
		result.CharacterExpectedDps[i] = agg.ToDescriptiveStats(b.characterExpectedDPS[i])
	}

	result.ExpectedDps = agg.ToDescriptiveStats(b.totalExpectedDPS)
}
