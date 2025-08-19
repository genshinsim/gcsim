package failures

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
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

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		var energy, stamina, swap, skill, dash, burstcd float64

		for _, fail := range result.Characters[i].FailedActions {
			switch fail.Reason {
			case action.InsufficientEnergy.String():
				energy += float64(fail.End-fail.Start) / 60
			case action.InsufficientStamina.String():
				stamina += float64(fail.End-fail.Start) / 60
			case action.SwapCD.String():
				swap += float64(fail.End-fail.Start) / 60
			case action.SkillCD.String():
				skill += float64(fail.End-fail.Start) / 60
			case action.DashCD.String():
				dash += float64(fail.End-fail.Start) / 60
			case action.BurstCD.String():
				burstcd += float64(fail.End-fail.Start) / 60
			}
		}

		b.failures[i].energy.Add(energy)
		b.failures[i].stamina.Add(stamina)
		b.failures[i].swap.Add(swap)
		b.failures[i].skill.Add(skill)
		b.failures[i].dash.Add(dash)
		b.failures[i].burstcd.Add(burstcd)
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
