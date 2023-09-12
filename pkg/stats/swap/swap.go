package swap

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	activeChar      int
	activeCharStart int
	activeIntervals []stats.ActiveCharacterInterval
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{
		activeChar:      core.Player.Active(),
		activeCharStart: 0,
	}

	core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		active := args[1].(int)

		interval := stats.ActiveCharacterInterval{
			Start:     out.activeCharStart,
			End:       core.F,
			Character: out.activeChar,
		}
		out.activeIntervals = append(out.activeIntervals, interval)
		out.activeChar = active
		out.activeCharStart = core.F

		return false
	}, "stats-swap-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	interval := stats.ActiveCharacterInterval{
		Start:     b.activeCharStart,
		End:       core.F,
		Character: b.activeChar,
	}
	result.ActiveCharacters = b.activeIntervals
	result.ActiveCharacters = append(result.ActiveCharacters, interval)
}
