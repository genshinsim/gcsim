package frostbearer

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("frostbearer", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

	icd := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if s.Frame() > icd {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		if s.Rand().Float64() < p {
			icd = s.Frame() + 600
			d := c.Snapshot(
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
			c.QueueDmg(&d, 1)

		}
	}, fmt.Sprintf("forstbearer-%v", c.Name()))
}
