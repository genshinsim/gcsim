package generic

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("generic bow", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	atk := 1.8 + float64(r)*0.6
	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if s.Frame() > icd {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if s.Rand().Float64() < 0.5 {
			icd = s.Frame() + 900 //15 sec icd
			d := c.Snapshot(
				"Prototype Archaic Proc",
				def.AttackTagWeaponSkill,
				def.ICDTagNone,
				def.ICDGroupDefault,
				def.StrikeTypeDefault,
				def.Physical,
				100,
				atk,
			)
			d.Targets = def.TargetAll
			c.QueueDmg(&d, 1)
		}
	}, fmt.Sprintf("forstbearer-%v", c.Name()))
}
