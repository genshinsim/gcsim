package kokomi

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Passive 2 - permanently modify stats for +25% healing bonus and -100% CR
func (c *char) passive() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = .25
	m[attributes.CR] = -1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("kokomi-passive", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("kokomiburst") == 0 {
			return false
		}

		a4Bonus := c.Stat(attributes.Heal) * 0.15 * c.MaxHP()
		atk.Info.FlatDmg += a4Bonus

		return false
	}, "kokomi-a4")
}
