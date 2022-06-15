package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func (c *char) c1() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		e, ok := args[0].(core.Enemy)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil != "Frost Flake Arrow" {
			return false
		}

		c.AddEnergy("ganyu-c1", 2)
		e.AddResistMod("ganyu-c1", 5*60, attributes.Cryo, -0.15)

		return false
	}, "ganyu-c1")
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod("ganyu-c4", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			x, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			// reset stacks on expiry
			if c.Core.F > x.GetTag("ganyuc4") {
				c.c4Stacks = 0
			}
			m[attributes.DmgP] = float64(c.c4Stacks) * 0.05
			return m, c.c4Stacks > 0
		})
	}
}
