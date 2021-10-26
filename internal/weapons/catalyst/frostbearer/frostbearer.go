package frostbearer

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("frostbearer", weapon)
	core.RegisterWeaponFunc("frostbearer", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F > icd {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Rand.Float64() < p {
			icd = c.F + 600
			d := char.Snapshot(
				"Frostbearer Proc",
				core.AttackTagWeaponSkill,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Physical,
				100,
				atk,
			)
			d.Targets = core.TargetAll
			if t.AuraType() == core.Cryo || t.AuraType() == core.Frozen {
				d.Mult = atkc
			}
			char.QueueDmg(&d, 1)

		}
		return false
	}, fmt.Sprintf("forstbearer-%v", char.Name()))
}
