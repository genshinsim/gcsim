package reaction

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var eventToReaction = map[event.Event]model.ReactionType{
	event.OnOverload:           model.ReactionTypeOverload,
	event.OnSuperconduct:       model.ReactionTypeSuperconduct,
	event.OnMelt:               model.ReactionTypeMelt,
	event.OnVaporize:           model.ReactionTypeVaporize,
	event.OnFrozen:             model.ReactionTypeFreeze,
	event.OnShatter:            model.ReactionTypeShatter,
	event.OnElectroCharged:     model.ReactionTypeElectroCharged,
	event.OnSwirlHydro:         model.ReactionTypeSwirlHydro,
	event.OnSwirlCryo:          model.ReactionTypeSwirlCryo,
	event.OnSwirlElectro:       model.ReactionTypeSwirlElectro,
	event.OnSwirlPyro:          model.ReactionTypeSwirlPyro,
	event.OnCrystallizeCryo:    model.ReactionTypeCrystallizeCryo,
	event.OnCrystallizeElectro: model.ReactionTypeCrystallizeElectro,
	event.OnCrystallizeHydro:   model.ReactionTypeCrystallizeHydro,
	event.OnCrystallizePyro:    model.ReactionTypeCrystallizePyro,
	event.OnAggravate:          model.ReactionTypeAggravate,
	event.OnSpread:             model.ReactionTypeSpread,
	event.OnQuicken:            model.ReactionTypeQuicken,
	event.OnBloom:              model.ReactionTypeBloom,
	event.OnHyperbloom:         model.ReactionTypeHyperbloom,
	event.OnBurgeon:            model.ReactionTypeBurgeon,
	event.OnBurning:            model.ReactionTypeBurning,
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

	eventSubFunc := func(reaction model.ReactionType) func(args ...interface{}) bool {
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
