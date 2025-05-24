package heal

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(stats.Config{
		Name: "heal",
		New:  NewStat,
	})
}

type buffer struct {
	events [][]stats.HealEvent
}

func NewStat(core *core.Core) (stats.Collector, error) {
	partySize := len(core.Player.Chars())
	out := buffer{
		events: make([][]stats.HealEvent, partySize),
	}

	core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		info := args[0].(*info.HealInfo)
		// Abyss card buff- heal may originate from no character at all. Do not log.
		if info.Caller < 0 || info.Caller >= partySize {
			return false
		}

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
