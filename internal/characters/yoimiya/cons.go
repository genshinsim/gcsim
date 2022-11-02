package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		trg, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}
		// ignore if debuff not on enemy
		if !trg.StatusIsActive(abDebuff) {
			return false
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("yoimiya-c1", 1200),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, "yoimiya-c1")
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.25
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)

		if atk.Info.ActorIndex != c.Index || !crit {
			return false
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("yoimiya-c2", 360),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, "yoimiya-c2")
}
