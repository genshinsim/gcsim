package expectedcritdmg

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// 30 = .5s
const bucketSize int = 30

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	expectedDmgCumu []float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	if !core.Flags.ExpectedCritDmg {
		out := buffer{}
		return &out, nil
	}
	out := buffer{
		expectedDmgCumu: make([]float64, len(core.Player.Chars())),
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
		out.expectedDmgCumu[attack.Info.ActorIndex] += damage * (1.0 + cr*cd)

		return false
	}, "stats-avg-crit-dmg-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	if !core.Flags.ExpectedCritDmg {
		return
	}
	for c, v := range b.expectedDmgCumu {
		result.Characters[c].ExpectedCritDmg = v
	}
}
