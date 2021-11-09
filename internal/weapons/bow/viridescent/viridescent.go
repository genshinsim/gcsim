package viridescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the viridescent hunt", weapon)
	core.RegisterWeaponFunc("theviridescenthunt", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	cd := 900 - r*60
	icd := 0
	mult := 0.3 + float64(r)*0.1

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		//check if char is correct?
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		// Vhunt passive only applies for NAs and CAs
		// For Tartaglia this also includes melee NAs/CAs
		// See https://youtu.be/EBtOiFhrs94?t=221, Test 4 and 5
		if !((ds.AttackTag == core.AttackTagNormal) || (ds.AttackTag == core.AttackTagExtra)) {
			return false
		}
		//check if cd is up
		if icd > c.F {
			return false
		}
		//50% chance to proc
		if c.Rand.Float64() > 0.5 {
			return false
		}

		//add a new action that deals % dmg immediately
		d := char.Snapshot(
			"Viridescent",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			mult,
		)
		d.Targets = core.TargetAll
		for i := 0; i <= 240; i += 30 {
			x := d.Clone()
			char.QueueDmg(&x, i)
		}

		//trigger cd
		icd = c.F + cd

		return false
	}, fmt.Sprintf("veridescent-%v", char.Name()))

}
