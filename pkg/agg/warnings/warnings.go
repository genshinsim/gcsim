package warnings

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
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
	burstcd calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		energy:  calc.StreamStats{},
		stamina: calc.StreamStats{},
		swap:    calc.StreamStats{},
		skill:   calc.StreamStats{},
		dash:    calc.StreamStats{},
		burstcd: calc.StreamStats{},
	}
	return &out, nil
}

func (b *buffer) Add(result *agg.Summary) {
	b.energy.Add(result.Failures.Energy)
	b.stamina.Add(result.Failures.Stamina)
	b.swap.Add(result.Failures.Swap)
	b.skill.Add(result.Failures.Skill)
	b.dash.Add(result.Failures.Dash)
	b.burstcd.Add(result.Failures.BurstCD)
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
		BurstCd:             b.burstcd.Mean() >= 1.0,
	}
}
