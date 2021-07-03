package freedom

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("freedom-sworn", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.DmgP] = 0.075 + float64(r)*0.025
	c.AddMod(def.CharStatMod{
		Key: "freedom-dmg",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = .15 + float64(r)*0.05
	plunge := .12 + 0.4*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnReaction(func(t def.Target, ds *def.Snapshot) {
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
		icd = s.Frame() + 30
		stacks++
		if stacks == 2 {
			stacks = 0
			s.AddStatus("freedom", 720)
			cooldown = s.Frame() + 1200
			for _, char := range s.Characters() {
				char.AddMod(def.CharStatMod{
					Key: "freedom-proc",
					Amount: func(a def.AttackTag) ([]float64, bool) {
						val[def.DmgP] = 0
						if a == def.AttackTagNormal || a == def.AttackTagExtra || a == def.AttackTagPlunge {
							val[def.DmgP] = plunge
						}
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}
	}, fmt.Sprintf("freedom-%v", c.Name()))

}
