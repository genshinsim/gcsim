package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.12

	icd := -1
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		// TODO: it works off field?
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !crit {
			return false
		}
		if icd > c.Core.F {
			return false
		}

		c.AddAttackMod(
			"xinyan-c1",
			5*60,
			func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return m, true
			},
		)
		icd = c.Core.F + 5*60

		return false
	}, "xinyan-c1")
}
