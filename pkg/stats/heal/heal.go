package heal

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	events [][]stats.HealEvent
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events: make([][]stats.HealEvent, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		info := args[0].(*player.HealInfo)
		target := args[1].(int)
		amount := args[2].(float64)

		event := stats.HealEvent{
			Frame:  core.F,
			Source: info.Message,
			Target: target,
			Heal:   amount,
		}
		out.events[info.Caller] = append(out.events[info.Caller], event)

		return false
	}, "stats-heal-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].HealEvents = b.events[c]
	}
}
