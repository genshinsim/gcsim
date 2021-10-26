package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward harp", weapon)
	core.RegisterWeaponFunc("skywardharp", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	//add passive crit, atk speed not sure how to do right now??
	//looks like jsut reduce the frames of normal attacks by 1 + 12%
	m := make([]float64, core.EndStatType)
	m[core.CD] = 0.15 + float64(r)*0.05
	cd := 270 - 30*r
	p := 0.5 + 0.1*float64(r)
	char.AddMod(core.CharStatMod{
		Key: "skyward harp",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		//check if char is correct?
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		//check if cd is up
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() > p {
			return false
		}

		//add a new action that deals % dmg immediately
		d := char.Snapshot(
			"Skyward Harp Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			1.25,
		)
		d.Targets = core.TargetAll
		char.QueueDmg(&d, 1)

		//trigger cd
		icd = c.F + cd

		return false
	}, fmt.Sprintf("skyward-harp-%v", char.Name()))

}
