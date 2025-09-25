package position

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(stats.Config{
		Name: "position",
		New:  NewStat,
	})
}

type buffer struct {
	targetOverlap bool
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		targetOverlap: overlaps(core.Combat.Enemies()),
	}

	core.Events.Subscribe(event.OnTargetMoved, func(args ...any) bool {
		target := args[0].(info.Target)

		for _, enemy := range core.Combat.Enemies() {
			if enemy.Key() == target.Key() {
				continue
			}
			if target.WillCollide(enemy.Shape()) {
				out.targetOverlap = true
			}
		}

		return false
	}, "stats-target-overlap")

	return &out, nil
}

func overlaps(targets []info.Target) bool {
	for i := range targets {
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
