package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype archaic", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := 1.8 + float64(r)*0.6
	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F > icd {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			icd = c.F + 900 //15 sec icd
			d := char.Snapshot(
				"Prototype Archaic Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				atk,
			)
			d.Targets = core.TargetAll
			char.QueueDmg(&d, 1)
		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))

}
