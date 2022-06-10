package yunjin

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// After Cliffbreaker's Banner is unleashed, all nearby party members' Normal Attack DMG is increased by 15% for 12s.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod("yunjin-c2", 12*60, func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if ae.Info.AttackTag == combat.AttackTagNormal {
				return m, true
			}
			return nil, false
		})
	}
}

// When Yun Jin triggers the Crystallize Reaction, her DEF is increased by 20% for 12s.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = .2
	charModFunc := func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}

		c.AddStatMod("yunjin-c4", 12*60, attributes.DEFP, func() ([]float64, bool) {
			return m, true
		})

		return false
	}
	c.Core.Events.Subscribe(event.OnCrystallizeCryo, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeHydro, charModFunc, "yunjin-c4")
}

// Characters under the effects of the Flying Cloud Flag Formation have their Normal ATK SPD increased by 12%.
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = .12
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod("yunjin-c6", 12*60, attributes.AtkSpd, func() ([]float64, bool) {
			return m, true
		})
	}
}
