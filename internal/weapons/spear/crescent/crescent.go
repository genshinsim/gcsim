package crescent

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("crescent pike", weapon)
}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := .15 + float64(r)*.05
	active := 0

	c.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		c.Log.Debugw("crescent pike active", "event", core.LogWeaponEvent, "frame", c.F, "char", char.CharIndex(), "expiry", c.F+300)
		active = c.F + 300

		return false
	}, fmt.Sprintf("cp-%v", char.Name()))

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		//check if char is correct?
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.F < active {
			//add a new action that deals % dmg immediately
			d := char.Snapshot(
				"Crescent Pike Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				atk,
			)
			char.QueueDmg(&d, 1)
		}
		return false
	}, fmt.Sprintf("cpp-%v", char.Name()))

}
