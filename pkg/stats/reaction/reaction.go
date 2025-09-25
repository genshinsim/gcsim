package reaction

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var eventToReaction = map[event.Event]info.ReactionType{
	event.OnOverload:           info.ReactionTypeOverload,
	event.OnSuperconduct:       info.ReactionTypeSuperconduct,
	event.OnMelt:               info.ReactionTypeMelt,
	event.OnVaporize:           info.ReactionTypeVaporize,
	event.OnFrozen:             info.ReactionTypeFreeze,
	event.OnShatter:            info.ReactionTypeShatter,
	event.OnElectroCharged:     info.ReactionTypeElectroCharged,
	event.OnSwirlHydro:         info.ReactionTypeSwirlHydro,
	event.OnSwirlCryo:          info.ReactionTypeSwirlCryo,
	event.OnSwirlElectro:       info.ReactionTypeSwirlElectro,
	event.OnSwirlPyro:          info.ReactionTypeSwirlPyro,
	event.OnCrystallizeCryo:    info.ReactionTypeCrystallizeCryo,
	event.OnCrystallizeElectro: info.ReactionTypeCrystallizeElectro,
	event.OnCrystallizeHydro:   info.ReactionTypeCrystallizeHydro,
	event.OnCrystallizePyro:    info.ReactionTypeCrystallizePyro,
	event.OnAggravate:          info.ReactionTypeAggravate,
	event.OnSpread:             info.ReactionTypeSpread,
	event.OnQuicken:            info.ReactionTypeQuicken,
	event.OnBloom:              info.ReactionTypeBloom,
	event.OnHyperbloom:         info.ReactionTypeHyperbloom,
	event.OnBurgeon:            info.ReactionTypeBurgeon,
	event.OnBurning:            info.ReactionTypeBurning,
}

func init() {
	stats.Register(stats.Config{
		Name: "reaction",
		New:  NewStat,
	})
}

type buffer struct {
	events [][]stats.ReactionEvent
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events: make([][]stats.ReactionEvent, len(core.Player.Chars())),
	}

	eventSubFunc := func(reaction info.ReactionType) func(args ...any) bool {
		return func(args ...any) bool {
			target := args[0].(info.Target)
			attack := args[1].(*info.AttackEvent)

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
