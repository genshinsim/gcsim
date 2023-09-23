package energy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	events [][]stats.EnergyEvent
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events: make([][]stats.EnergyEvent, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		character := args[0].(*character.CharWrapper)
		preEnergy := args[1].(float64)
		amount := args[2].(float64)
		source := args[3].(string)

		event := stats.EnergyEvent{
			Frame:   core.F,
			Source:  source,
			Gained:  character.Energy - preEnergy,
			Wasted:  preEnergy + amount - character.Energy,
			Current: character.Energy,
		}

		if core.Player.Active() == character.Index {
			event.FieldStatus = stats.OnField
		} else {
			event.FieldStatus = stats.OffField
		}

		out.events[character.Index] = append(out.events[character.Index], event)
		return false
	}, "stats-energy-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].EnergyEvents = b.events[c]
	}
}
