package pines

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("song of broken pines", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.12 + float64(r)*0.04
	c.AddMod(core.CharStatMod{
		Key: "pines-atk",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)
	val[core.AtkSpd] = 0.09 + 0.03*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
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
				char.AddMod(core.CharStatMod{
					Key: "pines-proc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}

	}, fmt.Sprintf("pines-%v", c.Name()))
}
