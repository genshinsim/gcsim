package generic

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("elegy of the end", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.EM] = 45 + float64(r)*15
	c.AddMod(def.CharStatMod{
		Key: "eledgy-em",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = .15 + float64(r)*0.05
	val[def.EM] = 75 + float64(r)*25

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagElementalArt && ds.AttackTag != def.AttackTagElementalBurst {
			return
		}
		if cooldown > s.Frame() {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 12
		stacks++
		if stacks == 4 {
			stacks = 0
			s.AddStatus("elegy", 720)
			cooldown = s.Frame() + 1200
			for _, char := range s.Characters() {
				char.AddMod(def.CharStatMod{
					Key: "eledgy-proc",
					Amount: func(a def.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}

	}, fmt.Sprintf("elegy-%v", c.Name()))

}
