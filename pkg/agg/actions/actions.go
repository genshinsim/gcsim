package actions

import (
	"math"

	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	failures []charFailures
}

type charFailures struct {
	energy  calc.StreamStats
	stamina calc.StreamStats
	swap    calc.StreamStats
	skill   calc.StreamStats
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		failures: make([]charFailures, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.failures[i] = charFailures{
			energy:  calc.StreamStats{},
			stamina: calc.StreamStats{},
			swap:    calc.StreamStats{},
			skill:   calc.StreamStats{},
		}
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i, c := range result.Characters {
		var energy, stamina, swap, skill float64

		for _, fail := range c.FailedActions {
			switch fail.Reason {
			case action.InsufficientEnergy.String():
				energy += float64(fail.End-fail.Start) / 60
			case action.InsufficientStamina.String():
				stamina += float64(fail.End-fail.Start) / 60
			case action.SwapCD.String():
				swap += float64(fail.End-fail.Start) / 60
			case action.SkillCD.String():
				skill += float64(fail.End-fail.Start) / 60
			}
		}

		b.failures[i].energy.Add(energy)
		b.failures[i].stamina.Add(stamina)
		b.failures[i].swap.Add(swap)
		b.failures[i].skill.Add(skill)
	}
}

func (b *buffer) Flush(result *agg.Result) {
	result.FailedActions = make([]agg.FailedActions, len(b.failures))
	for i, c := range b.failures {
		result.FailedActions[i] = agg.FailedActions{
			InsufficientEnergy:  toFloatStat(c.energy),
			InsufficientStamina: toFloatStat(c.stamina),
			SwapCD:              toFloatStat(c.swap),
			SkillCD:             toFloatStat(c.skill),
		}
	}
}

func toFloatStat(input calc.StreamStats) agg.FloatStat {
	out := agg.FloatStat{
		Min:  input.Min,
		Max:  input.Max,
		Mean: input.Mean(),
		SD:   input.StdDev(),
	}
	if math.IsNaN(out.SD) {
		out.SD = 0
	}
	return out
}
