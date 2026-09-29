package failures

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	agg.Register(agg.Config{
		Name: "failures",
		New:  NewAgg,
	})
}

type buffer struct {
	failures []charFailures
}

type charFailures struct {
	energy  *calc.StreamStats
	stamina *calc.StreamStats
	swap    *calc.StreamStats
	skill   *calc.StreamStats
	dash    *calc.StreamStats
	burstcd *calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		failures: make([]charFailures, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.failures[i] = charFailures{
			energy:  &calc.StreamStats{},
			stamina: &calc.StreamStats{},
			swap:    &calc.StreamStats{},
			skill:   &calc.StreamStats{},
			dash:    &calc.StreamStats{},
			burstcd: &calc.StreamStats{},
		}
	}

	return &out, nil
}

func (b *buffer) Add(result *agg.Summary) {
	for i := range result.Characters {
		fail := &result.Characters[i].Failures
		b.failures[i].energy.Add(fail.Energy)
		b.failures[i].stamina.Add(fail.Stamina)
		b.failures[i].swap.Add(fail.Swap)
		b.failures[i].skill.Add(fail.Skill)
		b.failures[i].dash.Add(fail.Dash)
		b.failures[i].burstcd.Add(fail.BurstCD)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.FailedActions = make([]*model.FailedActions, len(b.failures))
	for i, c := range b.failures {
		result.FailedActions[i] = &model.FailedActions{
			InsufficientEnergy:  agg.ToDescriptiveStats(c.energy),
			InsufficientStamina: agg.ToDescriptiveStats(c.stamina),
			SwapCd:              agg.ToDescriptiveStats(c.swap),
			SkillCd:             agg.ToDescriptiveStats(c.skill),
			DashCd:              agg.ToDescriptiveStats(c.dash),
			BurstCd:             agg.ToDescriptiveStats(c.burstcd),
		}
	}
}
