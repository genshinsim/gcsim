package damage

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// 30 = .5s
const bucketSize int = 30

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	events  [][]stats.DamageEvent
	buckets []float64
	cumu    [][]float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events:  make([][]stats.DamageEvent, len(core.Player.Chars())),
		buckets: make([]float64, 0),
		cumu:    make([][]float64, 0),
	}
	out.cumu = append(out.cumu, make([]float64, len(core.Player.Chars())))

	core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		attack := args[1].(*combat.AttackEvent)
		damage := args[2].(float64)
		crit := args[3].(bool)

		// TODO: validate if this is still true?
		// No need to pull damage stats for non-enemies
		if target.Type() != targets.TargettableEnemy {
			return false
		}

		bucket := core.F / bucketSize
		last := out.cumu[len(out.cumu)-1]
		for bucket >= len(out.cumu) {
			newBucket := make([]float64, len(core.Player.Chars()))
			copy(newBucket, last)
			out.cumu = append(out.cumu, newBucket)
		}
		out.cumu[bucket][attack.Info.ActorIndex] += damage

		for bucket >= len(out.buckets) {
			out.buckets = append(out.buckets, float64(0))
		}
		out.buckets[bucket] += damage

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
			switch attack.Info.AmpType {
			case reactions.Vaporize:
				event.ReactionModifier = stats.Vaporize
			case reactions.Melt:
				event.ReactionModifier = stats.Melt
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
	result.DamageBuckets = b.buckets
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].DamageEvents = b.events[c]
		result.Characters[c].DamageCumulativeContrib = make([]float64, len(b.buckets))
	}

	for i := 0; i < len(b.cumu); i++ {
		var total float64
		for _, v := range b.cumu[i] {
			total += v
		}

		if total > 0 {
			for c, v := range b.cumu[i] {
				result.Characters[c].DamageCumulativeContrib[i] = v / total
			}
		}
	}
}
