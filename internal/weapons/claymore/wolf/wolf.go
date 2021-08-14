package generic

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("wolf's gravestone", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)
	c.AddMod(core.CharStatMod{
		Key:    "wolf-flat",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

	bonus := make([]float64, core.EndStatType)
	bonus[core.ATKP] = 0.3 + 0.1*float64(r)
	icd := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}
		if !s.Flags().DamageMode {
			return //ignore as we not tracking HP
		}
		if t.HP()/t.MaxHP() > 0.3 {
			return
		}
		icd = s.Frame() + 1800 //every 30 seconds

		for _, char := range s.Characters() {
			char.AddMod(core.CharStatMod{
				Key:    "wolf-proc",
				Expiry: s.Frame() + 720,
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return bonus, true
				},
			})
		}

	}, fmt.Sprintf("wolf-%v", c.Name()))
}
