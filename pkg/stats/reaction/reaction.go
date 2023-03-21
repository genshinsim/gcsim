package reaction

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var eventToReaction = map[event.Event]reactions.ReactionType{
	event.OnOverload:           reactions.Overload,
	event.OnSuperconduct:       reactions.Superconduct,
	event.OnMelt:               reactions.Melt,
	event.OnVaporize:           reactions.Vaporize,
	event.OnFrozen:             reactions.Freeze,
	event.OnElectroCharged:     reactions.ElectroCharged,
	event.OnSwirlHydro:         reactions.SwirlHydro,
	event.OnSwirlCryo:          reactions.SwirlCryo,
	event.OnSwirlElectro:       reactions.SwirlElectro,
	event.OnSwirlPyro:          reactions.SwirlPyro,
	event.OnCrystallizeCryo:    reactions.CrystallizeCryo,
	event.OnCrystallizeElectro: reactions.CrystallizeElectro,
	event.OnCrystallizeHydro:   reactions.CrystallizeHydro,
	event.OnCrystallizePyro:    reactions.CrystallizePyro,
	event.OnAggravate:          reactions.Aggravate,
	event.OnSpread:             reactions.Spread,
	event.OnQuicken:            reactions.Quicken,
	event.OnBloom:              reactions.Bloom,
	event.OnHyperbloom:         reactions.Hyperbloom,
	event.OnBurgeon:            reactions.Burgeon,
	event.OnBurning:            reactions.Burning,
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

	eventSubFunc := func(reaction reactions.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			target := args[0].(combat.Target)
			attack := args[1].(*combat.AttackEvent)

			event := stats.ReactionEvent{
				Frame:    attack.SourceFrame,
				Source:   attack.Info.Abil,
				Target:   int(target.Key()),
				Reaction: string(reaction),
			}
			out.events[attack.Info.ActorIndex] = append(out.events[attack.Info.ActorIndex], event)
			return false
		}
	}

	for k, v := range eventToReaction {
		core.Events.Subscribe(k, eventSubFunc(v), "stats-reaction-log")
	}

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].ReactionEvents = b.events[c]
	}
}
