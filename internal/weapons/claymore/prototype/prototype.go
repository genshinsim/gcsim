package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("prototype archaic", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	atk := 1.8 + float64(r)*0.6
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
		if s.Rand().Float64() < 0.5 {
			icd = s.Frame() + 900 //15 sec icd
			d := c.Snapshot(
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
			c.QueueDmg(&d, 1)
		}
	}, fmt.Sprintf("forstbearer-%v", c.Name()))
}
