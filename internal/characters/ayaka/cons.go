package ayaka

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// Callback for Ayaka C1 that is attached to NA/CA hits
// When Kamisato Ayaka's Normal or Charged Attacks deal Cryo DMG to opponents, it has a 50% chance of decreasing the CD of Kamisato Art: Hyouka by 0.3s.
// This effect can occur once every 0.1s.
func (c *char) c1(a combat.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}

	if a.AttackEvent.Info.Element != attributes.Cryo {
		return
	}
	if c.icdC1 > c.Core.F {
		return
	}
	if c.Core.Rand.Float64() < .5 {
		return
	}
	c.ReduceActionCooldown(action.ActionSkill, 18)
	c.icdC1 = c.Core.F + 6
}

// Callback for Ayaka C4 that is attached to Burst hits
// Opponents damaged by Kamisato Art: Soumetsu's Frostflake Seki no To will have their DEF decreased by 30% for 6s.
func (c *char) c4(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}

	e, ok := a.Target.(core.Enemy)
	if !ok {
		return
	}
	e.AddDefMod("ayaka-c4", 60*6, -0.3)
}

// Callback for Ayaka C6 that is attached to CA hits
func (c *char) c6(a combat.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}

	if !c.c6CDTimerAvail {
		return
	}
	c.c6CDTimerAvail = false

	c.Core.Tasks.Add(func() {
		c.DeleteAttackMod("ayaka-c6")

		c.Core.Tasks.Add(func() {
			c.c6CDTimerAvail = true
			c.c6AddBuff()
		}, 600)
	}, 30)
}

func (c *char) c6AddBuff() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 2.98

	c.AddAttackMod("ayaka-c6", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag != combat.AttackTagExtra {
			return nil, false
		}
		return m, true
	})
}
