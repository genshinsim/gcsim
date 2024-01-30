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
	stats.Register(stats.Config{
		Name: "damage",
		New:  NewStat,
	})
}

type buffer struct {
	events     [][]stats.DamageEvent
	buckets    []float64
	cumuChar   [][]float64
	cumuTarget [][]float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events:     make([][]stats.DamageEvent, len(core.Player.Chars())),
		buckets:    make([]float64, 0),
		cumuChar:   make([][]float64, 0),
		cumuTarget: make([][]float64, 0),
	}
	out.cumuChar = append(out.cumuChar, make([]float64, len(core.Player.Chars())))
	out.cumuTarget = append(out.cumuTarget, make([]float64, len(core.Combat.Enemies())))

	core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		targetKey := target.Key()
		attack := args[1].(*combat.AttackEvent)
		damage := args[2].(float64)
		crit := args[3].(bool)

		// TODO: validate if this is still true?
		// No need to pull damage stats for non-enemies
		if target.Type() != targets.TargettableEnemy {
			return false
		}

		bucket := core.F / bucketSize

		last := out.cumuChar[len(out.cumuChar)-1]
		for bucket >= len(out.cumuChar) {
			newBucket := make([]float64, len(core.Player.Chars()))
			copy(newBucket, last)
			out.cumuChar = append(out.cumuChar, newBucket)
		}
		out.cumuChar[bucket][attack.Info.ActorIndex] += damage

		last = out.cumuTarget[len(out.cumuTarget)-1]
		for bucket >= len(out.cumuTarget) {
			newBucket := make([]float64, len(core.Combat.Enemies()))
			copy(newBucket, last)
			out.cumuTarget = append(out.cumuTarget, newBucket)
		}
		// TODO: subject to break if target key gen changes...
		out.cumuTarget[bucket][targetKey-1] += damage

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
			Target:  int(targetKey),
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

	for i := 0; i < len(b.cumuChar); i++ {
		var total float64
		for _, v := range b.cumuChar[i] {
			total += v
		}

		if total > 0 {
			for c, v := range b.cumuChar[i] {
				result.Characters[c].DamageCumulativeContrib[i] = v / total
			}
		}
	}

	// TODO: working under assumption that enemies are not removed from handler array upon death, subject to break...
	bucketCount := len(result.DamageBuckets)
	for e := range core.Combat.Enemies() {
		result.Enemies[e].CumulativeDamage = make([]float64, bucketCount)
	}
	if bucketCount == 0 {
		return
	}
	for i := 0; i < len(b.cumuTarget); i++ {
		for e, v := range b.cumuTarget[i] {
			result.Enemies[e].CumulativeDamage[i] = v
		}
	}
}
