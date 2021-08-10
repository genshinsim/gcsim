package elegy

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("elegy of the end", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.EM] = 45 + float64(r)*15
	c.AddMod(core.CharStatMod{
		Key: "eledgy-em",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = .15 + float64(r)*0.05
	val[core.EM] = 75 + float64(r)*25

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagElementalArt && ds.AttackTag != core.AttackTagElementalBurst {
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
				char.AddMod(core.CharStatMod{
					Key: "eledgy-proc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}

	}, fmt.Sprintf("elegy-%v", c.Name()))

}
