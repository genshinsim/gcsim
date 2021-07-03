package starsilver

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("snow-tombed starsilver", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	atk := 0.65 + float64(r)*0.15
	atkc := 1.6 + float64(r)*0.4
	p := 0.5 + float64(r)*0.1

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
		if s.Rand().Float64() < p {
			icd = s.Frame() + 600
			d := c.Snapshot(
				"Starsilver Proc",
				def.AttackTagWeaponSkill,
				def.ICDTagNone,
				def.ICDGroupDefault,
				def.StrikeTypeDefault,
				def.Physical,
				100,
				atk,
			)
			d.Targets = def.TargetAll
			if t.AuraType() == def.Cryo || t.AuraType() == def.Frozen {
				d.Mult = atkc
			}
			c.QueueDmg(&d, 1)

		}
	}, fmt.Sprintf("starsilver-%v", c.Name()))
}
