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

// When Yoimiya's Pyro DMG scores a CRIT Hit, Yoimiya will gain a 25% Pyro DMG Bonus for 6s.
// This effect can be triggered even when Yoimiya is not the active character.
func (c *char) makeC2CB() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if !a.IsCrit {
			return
		}
		if a.AttackEvent.Info.Element != attributes.Pyro {
			return
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.PyroP] = 0.25
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("yoimiya-c2", 360),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
