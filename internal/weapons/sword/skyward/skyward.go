package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward blade", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	dur := -1
	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		dur = c.F + 720
		c.Log.Debugw("Skyward Blade activated", "frame", c.F, "event", core.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.03 + float64(r)*0.01

	char.AddMod(core.CharStatMod{
		Key: "skyward blade",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.AtkSpd] = 0
			if dur > c.F {
				m[core.AtkSpd] = 0.1 //if burst active
			}
			return m, true
		},
		Expiry: -1,
	})

	//damage procs
	atk := .15 + .05*float64(r)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		ds := args[1].(*core.Snapshot)

		//check if char is correct?
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		//check if buff up
		if dur < c.F {
			return false
		}

		//add a new action that deals % dmg immediately
		d := char.Snapshot(
			"Skyward Blade Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		char.QueueDmg(&d, 1)
		return false

	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

}
