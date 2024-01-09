package optimization

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type CustomDamageStatsBuffer struct {
	ExpectedDmgCumu []float64
	duration        int
}

func OptimizerDmgStat(core *core.Core) (stats.CollectorCustomStats[CustomDamageStatsBuffer], error) {
	if !core.Flags.IgnoreBurstEnergy {
		// This data doesn't mean much without the IgnoreBurstEnergy flag set
		// So the stat collector disables itself when this flag isn't set
		return &CustomDamageStatsBuffer{}, nil
	}

	out := CustomDamageStatsBuffer{
		ExpectedDmgCumu: make([]float64, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		attack := args[1].(*combat.AttackEvent)
		damage := args[2].(float64)
		crit := args[3].(bool)

		// TODO: validate if this is still true?
		// No need to pull damage stats for non-enemies
		if target.Type() != targets.TargettableEnemy {
			return false
		}
		cr := attack.Snapshot.Stats[attributes.CR]
		cd := attack.Snapshot.Stats[attributes.CD]

		// TODO: Do we need to ensure that 1 + cd > 0 ?
		if crit && 1+cd > 0 {
			damage /= (1 + cd)
		}
		out.ExpectedDmgCumu[attack.Info.ActorIndex] += damage * (1.0 + cr*cd)

		return false
	}, "substat-opt-dmg-log")
	return &out, nil
}

func (b CustomDamageStatsBuffer) Flush(core *core.Core) CustomDamageStatsBuffer {
	b.duration = core.F
	return b
}

type CustomDamageAggBuffer struct {
	ExpectedDps     []float64
	CharExpectedDps [][]float64
}

func NewDamageAggBuffer(cfg *info.ActionList) CustomDamageAggBuffer {
	character_count := len(cfg.Characters)
	return CustomDamageAggBuffer{
		CharExpectedDps: make([][]float64, character_count),
	}
}

func (agg *CustomDamageAggBuffer) Add(b CustomDamageStatsBuffer) {
	char_count := len(b.ExpectedDmgCumu)
	totalExpectedDPS := 0.0
	for i := 0; i < char_count; i++ {
		charExpectedDps := b.ExpectedDmgCumu[i] / (float64(b.duration) / 60.0)
		agg.CharExpectedDps[i] = append(agg.CharExpectedDps[i], charExpectedDps)
		totalExpectedDPS += charExpectedDps
	}
	agg.ExpectedDps = append(agg.ExpectedDps, totalExpectedDPS)
}

func (agg *CustomDamageAggBuffer) Flush() {
	char_count := len(agg.CharExpectedDps)
	for i := 0; i < char_count; i++ {
		slices.Sort(agg.CharExpectedDps[i])
	}
	slices.Sort(agg.ExpectedDps)
}
