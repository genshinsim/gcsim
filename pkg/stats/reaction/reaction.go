package reaction

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var reactions = map[event.Event]combat.ReactionType{
	event.OnOverload:           combat.Overload,
	event.OnSuperconduct:       combat.Superconduct,
	event.OnMelt:               combat.Melt,
	event.OnVaporize:           combat.Vaporize,
	event.OnFrozen:             combat.Freeze,
	event.OnElectroCharged:     combat.ElectroCharged,
	event.OnSwirlHydro:         combat.SwirlHydro,
	event.OnSwirlCryo:          combat.SwirlCryo,
	event.OnSwirlElectro:       combat.SwirlElectro,
	event.OnSwirlPyro:          combat.SwirlPyro,
	event.OnCrystallizeCryo:    combat.CrystallizeCryo,
	event.OnCrystallizeElectro: combat.CrystallizeElectro,
	event.OnCrystallizeHydro:   combat.CrystallizeHydro,
	event.OnCrystallizePyro:    combat.CrystallizePyro,
	event.OnAggravate:          combat.Aggravate,
	event.OnSpread:             combat.Spread,
	event.OnQuicken:            combat.Quicken,
	event.OnBloom:              combat.Bloom,
	event.OnHyperbloom:         combat.Hyperbloom,
	event.OnBurgeon:            combat.Burgeon,
	event.OnBurning:            combat.Burning,
}

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	events [][]stats.ReactionEvent
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{
		events: make([][]stats.ReactionEvent, len(core.Player.Chars())),
	}

	eventSubFunc := func(reaction combat.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			target := args[0].(combat.Target)
			attack := args[1].(*combat.AttackEvent)

			event := stats.ReactionEvent{
				Frame:    attack.SourceFrame,
				Source:   attack.Info.Abil,
				Target:   target.Index(),
				Reaction: string(reaction),
			}
			out.events[attack.Info.ActorIndex] = append(out.events[attack.Info.ActorIndex], event)
			return false
		}
	}

	for k, v := range reactions {
		core.Events.Subscribe(k, eventSubFunc(v), "stats-reaction-log")
	}

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].ReactionEvents = b.events[c]
	}
}
