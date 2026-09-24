package status

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(stats.Config{
		Name: "status",
		New:  NewStat,
	})
}

type buffer struct {
	reactionUptime  []map[string]int
	enemyReactions  [][]stats.ReactionStatusInterval
	activeReactions []map[info.ReactionModKey]int
}

func enemyIndex(core *core.Core, e *enemy.Enemy) int {
	for i, enemyTarget := range core.Combat.Enemies() {
		if enemyTarget.Key() == e.Key() {
			return i
		}
	}
	return -1
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		reactionUptime:  make([]map[string]int, len(core.Combat.Enemies())),
		enemyReactions:  make([][]stats.ReactionStatusInterval, len(core.Combat.Enemies())),
		activeReactions: make([]map[info.ReactionModKey]int, len(core.Combat.Enemies())),
	}

	for i := 0; i < len(core.Combat.Enemies()); i++ {
		out.reactionUptime[i] = make(map[string]int)
		out.activeReactions[i] = make(map[info.ReactionModKey]int)
	}

	addInterval := func(enemyIndex int, key info.ReactionModKey, start, end int) {
		interval := stats.ReactionStatusInterval{
			Start: start,
			End:   end,
			Type:  key.String(),
		}
		out.enemyReactions[enemyIndex] = append(out.enemyReactions[enemyIndex], interval)
		out.reactionUptime[enemyIndex][key.String()] += end - start
	}

	core.Events.Subscribe(event.OnAuraDurabilityAdded, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		key, ok := args[1].(info.ReactionModKey)
		if !ok {
			return
		}
		idx := enemyIndex(core, e)
		if idx < 0 {
			return
		}
		if _, ok := out.activeReactions[idx][key]; !ok {
			out.activeReactions[idx][key] = core.F
		}
	}, "stats-status-added")

	core.Events.Subscribe(event.OnAuraDurabilityDepleted, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		idx := enemyIndex(core, e)
		if idx < 0 {
			return
		}

		key, ok := args[1].(info.ReactionModKey)
		if !ok {
			return
		}

		start, ok := out.activeReactions[idx][key]
		if !ok {
			return
		}
		addInterval(idx, key, start, core.F)
		delete(out.activeReactions[idx], key)
	}, "stats-status-aura-depleted")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for e := 0; e < len(core.Combat.Enemies()); e++ {
		for k, start := range b.activeReactions[e] {
			interval := stats.ReactionStatusInterval{
				Start: start,
				End:   core.F,
				Type:  k.String(),
			}
			b.enemyReactions[e] = append(b.enemyReactions[e], interval)
			b.reactionUptime[e][k.String()] += core.F - start
		}

		result.Enemies[e].ReactionStatus = b.enemyReactions[e]
		result.Enemies[e].ReactionUptime = b.reactionUptime[e]
	}
}
