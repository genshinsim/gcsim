package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		e, ok := args[0].(*enemy.Enemy)
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
		e.AddResistMod(enemy.ResistMod{Base: modifier.NewBase("ganyu-c1", 5*60), Ele: attributes.Cryo, Value: -0.15})

		return false
	}, "ganyu-c1")
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{Base: modifier.NewBase("ganyu-c4", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
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
		}})
	}
}
