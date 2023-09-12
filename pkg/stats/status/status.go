package status

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// 6 = .1s. TODO: figure out best bucket size
const bucketSize int = 1

func init() {
	stats.Register(NewStat)
}

/*
 * TODO: minimize using this as collector as much as possible (runs every frame). Stuff to maybe move
	 * Health -> OnPlayerDamage + OnHeal (health intervals)?
	 * damage mitigation -> intervals?
	 * energy -> OnEnergyChange (might not be worth?)
	 * reactions -> OnAuraDurabilityAdded + OnAuraDurabilityDepleted (need to look into if this works)
	* TODO: Add a "None" interval for enemy auras
*/
type buffer struct {
	maxEnemyLvl      int
	damageMitigation []float64

	activeTime []int
	charEnergy [][]float64
	charHealth [][]float64

	reactionUptime  []map[string]int
	enemyReactions  [][]stats.ReactionStatusInterval
	activeReactions []map[reactable.Modifier]int
}

func maxUpdate(arr []float64, index int, value float64) []float64 {
	for index >= len(arr) {
		return append(arr, value)
	}

	if value > arr[index] {
		arr[index] = value
	}
	return arr
}

func avgUpdate(arr []float64, index int, value float64) []float64 {
	for index >= len(arr) {
		arr = append(arr, float64(0))
	}

	arr[index] += value / float64(bucketSize)
	return arr
}

func damageMod(c *character.CharWrapper, elvl int) float64 {
	def := c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)
	defmod := 1 - (def / (def + 5*float64(elvl) + 500))
	// TODO: implement resmod
	resmod := 1.0

	return defmod * resmod
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		activeTime: make([]int, len(core.Player.Chars())),
		charEnergy: make([][]float64, len(core.Player.Chars())),
		charHealth: make([][]float64, len(core.Player.Chars())),

		reactionUptime:  make([]map[string]int, len(core.Combat.Enemies())),
		enemyReactions:  make([][]stats.ReactionStatusInterval, len(core.Combat.Enemies())),
		activeReactions: make([]map[reactable.Modifier]int, len(core.Combat.Enemies())),
	}

	for i := 0; i < len(core.Combat.Enemies()); i++ {
		out.reactionUptime[i] = make(map[string]int)
		out.activeReactions[i] = make(map[reactable.Modifier]int)

		if enemy, ok := core.Combat.Enemies()[i].(*enemy.Enemy); ok {
			if enemy.Level > out.maxEnemyLvl {
				out.maxEnemyLvl = enemy.Level
			}
		}
	}

	core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		bucket := int(core.F / bucketSize)
		active := core.Player.ActiveChar()

		out.activeTime[active.Index] += 1
		out.damageMitigation = avgUpdate(
			out.damageMitigation, bucket, damageMod(active, out.maxEnemyLvl))

		for i, char := range core.Player.Chars() {
			out.charEnergy[i] = maxUpdate(out.charEnergy[i], bucket, char.Energy)
			out.charHealth[i] = avgUpdate(out.charHealth[i], bucket, char.CurrentHP())
		}

		for i, t := range core.Combat.Enemies() {
			enemy, ok := t.(*enemy.Enemy)
			if !ok {
				continue
			}

			current := make(map[reactable.Modifier]int)

			for r, v := range enemy.Durability {
				if v <= reactable.ZeroDur {
					continue
				}
				var key = reactable.Modifier(r)
				out.reactionUptime[i][key.String()] += 1

				if start, ok := out.activeReactions[i][key]; ok {
					current[key] = start
				} else {
					current[key] = core.F
				}
			}

			for k, start := range out.activeReactions[i] {
				_, ok := current[k]
				if ok {
					continue
				}

				if core.F-start <= 5 {
					continue
				}

				interval := stats.ReactionStatusInterval{
					Start: start,
					End:   core.F,
					Type:  k.String(),
				}
				out.enemyReactions[i] = append(out.enemyReactions[i], interval)
			}

			out.activeReactions[i] = current
		}

		return false
	}, "stats-status-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	fill := bucketSize - (core.F % bucketSize) - 1
	bucket := int(core.F / bucketSize)

	// for averages, last bucket is inaccurate. Fill to fix
	for i := 0; i < fill; i++ {
		b.damageMitigation = avgUpdate(
			b.damageMitigation, bucket, damageMod(core.Player.ActiveChar(), b.maxEnemyLvl))
	}

	result.DamageMitigation = b.damageMitigation

	for c := 0; c < len(core.Player.Chars()); c++ {
		for i := 0; i < fill; i++ {
			b.charHealth[c] = avgUpdate(b.charHealth[c], bucket, core.Player.Chars()[c].CurrentHP())
		}

		result.Characters[c].ActiveTime = b.activeTime[c]
		result.Characters[c].HealthStatus = b.charHealth[c]
		result.Characters[c].EnergyStatus = b.charEnergy[c]
	}

	for e := 0; e < len(core.Combat.Enemies()); e++ {
		for k, start := range b.activeReactions[e] {
			if core.F-start <= 5 {
				continue
			}

			interval := stats.ReactionStatusInterval{
				Start: start,
				End:   core.F,
				Type:  k.String(),
			}
			b.enemyReactions[e] = append(b.enemyReactions[e], interval)
		}

		result.Enemies[e].ReactionStatus = b.enemyReactions[e]
		result.Enemies[e].ReactionUptime = b.reactionUptime[e]
	}
}
