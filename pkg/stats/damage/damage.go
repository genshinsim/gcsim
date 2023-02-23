package damage

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	events [][]stats.DamageEvent
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{
		events: make([][]stats.DamageEvent, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		attack := args[1].(*combat.AttackEvent)
		damage := args[2].(float64)
		crit := args[3].(bool)

		// TODO: validate if this is still true?
		// No need to pull damage stats for non-enemies
		if target.Type() != combat.TargettableEnemy {
			return false
		}

		// TODO: ActionId population
		// TODO: Modifiers population
		// TODO: Mitigation population
		event := stats.DamageEvent{
			Frame:   attack.SourceFrame,
			Source:  attack.Info.Abil,
			Target:  int(target.Key()),
			Element: attack.Info.Element.String(),
			Crit:    crit,
			Damage:  damage,
		}

		if attack.Info.Amped {
			switch attack.Info.AmpMult {
			case 1.5:
				event.ReactionModifier = stats.Amp15
			case 2:
				event.ReactionModifier = stats.Amp20
			}
		}

		if attack.Info.Catalyzed {
			switch attack.Info.CatalyzedType {
			case reactions.Aggravate:
				event.ReactionModifier = stats.Aggravate
			case reactions.Spread:
				event.ReactionModifier = stats.Spread
			}
		}

		out.events[attack.Info.ActorIndex] = append(out.events[attack.Info.ActorIndex], event)
		return false
	}, "stats-dmg-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].DamageEvents = b.events[c]
	}
}
