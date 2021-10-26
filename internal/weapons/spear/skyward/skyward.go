package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward spine", weapon)
	core.RegisterWeaponFunc("skywardspine", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.06 + float64(r)*0.02
	m[core.AtkSpd] = 0.12

	char.AddMod(core.CharStatMod{
		Key: "skyward spine",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0
	atk := .25 + .15*float64(r)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		//check if char is correct?
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		//check if cd is up
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() > .5 {
			return false
		}

		//add a new action that deals % dmg immediately
		d := char.Snapshot(
			"Skyward Spine Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		char.QueueDmg(&d, 1)

		//trigger cd
		icd = c.F + 120
		return false
	}, fmt.Sprintf("skyward-spine-%v", char.Name()))

}
