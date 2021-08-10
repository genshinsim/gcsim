package freedom

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("freedom-sworn", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.075 + float64(r)*0.025
	c.AddMod(core.CharStatMod{
		Key: "freedom-dmg",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = .15 + float64(r)*0.05
	plunge := .12 + 0.4*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	s.AddOnReaction(func(t core.Target, ds *core.Snapshot) {
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
		icd = s.Frame() + 30
		stacks++
		if stacks == 2 {
			stacks = 0
			s.AddStatus("freedom", 720)
			cooldown = s.Frame() + 1200
			for _, char := range s.Characters() {
				char.AddMod(core.CharStatMod{
					Key: "freedom-proc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						val[core.DmgP] = 0
						if a == core.AttackTagNormal || a == core.AttackTagExtra || a == core.AttackTagPlunge {
							val[core.DmgP] = plunge
						}
						return val, true
					},
					Expiry: s.Frame() + 720,
				})
			}
		}
	}, fmt.Sprintf("freedom-%v", c.Name()))

}
