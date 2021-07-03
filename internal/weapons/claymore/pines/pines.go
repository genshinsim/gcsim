package pines

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("song of broken pines", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.ATKP] = 0.12 + float64(r)*0.04
	c.AddMod(def.CharStatMod{
		Key: "pines-atk",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = 0.15 + 0.05*float64(r)
	val[def.AtkSpd] = 0.09 + 0.03*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
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
			s.AddStatus("pines", 720)
			cooldown = s.Frame() + 1200
			for _, char := range s.Characters() {
				char.AddMod(def.CharStatMod{
					Key: "pines-proc",
					Amount: func(a def.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}

	}, fmt.Sprintf("pines-%v", c.Name()))
}
