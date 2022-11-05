package position

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	targetOverlap bool
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{
		targetOverlap: overlaps(core.Combat.Enemies()),
	}

	core.Events.Subscribe(event.OnTargetMoved, func(args ...interface{}) bool {
		target := args[0].(combat.Target)

		for _, enemy := range core.Combat.Enemies() {
			if target.WillCollide(enemy.Shape()) {
				out.targetOverlap = true
			}
		}

		return false
	}, "stats-target-overlap")

	return &out, nil
}

func overlaps(targets []combat.Target) bool {
	for i := 0; i < len(targets); i++ {
		for j := i + 1; j < len(targets); j++ {
			if targets[i].WillCollide(targets[j].Shape()) {
				return true
			}
		}
	}
	return false
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	result.TargetOverlap = b.targetOverlap
}
