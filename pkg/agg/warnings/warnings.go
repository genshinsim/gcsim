package warnings

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
		Name: "warnings",
		New:  NewAgg,
	})
}

type buffer struct {
	overlap bool
	energy  calc.StreamStats
	stamina calc.StreamStats
	swap    calc.StreamStats
	skill   calc.StreamStats
	dash    calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		energy:  calc.StreamStats{},
		stamina: calc.StreamStats{},
		swap:    calc.StreamStats{},
		skill:   calc.StreamStats{},
		dash:    calc.StreamStats{},
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	var energy, stamina, swap, skill, dash float64

	for i := range result.Characters {
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
			}
		}
	}

	b.energy.Add(energy)
	b.stamina.Add(stamina)
	b.swap.Add(swap)
	b.skill.Add(skill)
	b.dash.Add(dash)
	b.overlap = b.overlap || result.TargetOverlap
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.Warnings = &model.Warnings{
		TargetOverlap:       b.overlap,
		InsufficientEnergy:  b.energy.Mean() >= 1.0,
		InsufficientStamina: b.stamina.Mean() >= 1.0,
		SwapCd:              b.swap.Mean() >= 1.0,
		SkillCd:             b.skill.Mean() >= 1.0,
		DashCd:              b.dash.Mean() >= 1.0,
	}
}
