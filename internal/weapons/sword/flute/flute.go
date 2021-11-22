package flute

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the flute", weapon)
	core.RegisterWeaponFunc("theflute", weapon)
}

//Normal or Charged Attacks grant a Harmonic on hits. Gaining 5 Harmonics triggers the
//power of music and deals 100% ATK DMG to surrounding opponents. Harmonics last up to 30s,
//and a maximum of 1 can be gained every 0.5s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	expiry := 0
	stacks := 0
	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		ds := args[1].(*core.Snapshot)

		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 30 // every .5 sec
		if expiry < c.F {
			stacks = 0
		}
		stacks++
		expiry = c.F + 1800 //stacks lasts 30s

		if stacks == 5 {
			//trigger dmg at 5 stacks
			stacks = 0
			expiry = 0

			d := char.Snapshot(
				"Flute Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				0.75+0.25*float64(r),
			)
			d.Targets = core.TargetAll
			char.QueueDmg(&d, 1)

		}
		return false
	}, fmt.Sprintf("prototype-rancour-%v", char.Name()))

}
