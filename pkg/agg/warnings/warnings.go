package warnings

import (
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
	overlap bool
	energy  calc.StreamStats
	stamina calc.StreamStats
	swap    calc.StreamStats
	skill   calc.StreamStats
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		energy:  calc.StreamStats{},
		stamina: calc.StreamStats{},
		swap:    calc.StreamStats{},
		skill:   calc.StreamStats{},
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	var energy, stamina, swap, skill float64

	for _, c := range result.Characters {
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
	}

	b.energy.Add(energy)
	b.stamina.Add(stamina)
	b.swap.Add(swap)
	b.skill.Add(skill)
}

func (b *buffer) Flush(result *agg.Result) {
	result.Warnings = agg.Warnings{
		TargetOverlap:       false, // TODO
		InsufficientEnergy:  b.energy.StdDev() >= 1.0,
		InsufficientStamina: b.stamina.Mean() >= 1.0,
		SwapCD:              b.swap.Mean() >= 1.0,
		SkillCD:             b.skill.Mean() >= 1.0,
	}
}
