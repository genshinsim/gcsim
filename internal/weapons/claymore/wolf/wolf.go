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

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = 0.15 + 0.05*float64(r)
	c.AddMod(def.CharStatMod{
		Key:    "wolf-flat",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

	bonus := make([]float64, def.EndStatType)
	bonus[def.ATKP] = 0.3 + 0.1*float64(r)
	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}
		if !s.Flags().HPMode {
			return //ignore as we not tracking HP
		}
		if t.HP()/t.MaxHP() > 0.3 {
			return
		}
		icd = s.Frame() + 1800 //every 30 seconds

		for _, char := range s.Characters() {
			char.AddMod(def.CharStatMod{
				Key:    "wolf-proc",
				Expiry: s.Frame() + 720,
				Amount: func(a def.AttackTag) ([]float64, bool) {
					return bonus, true
				},
			})
		}

	}, fmt.Sprintf("wolf-%v", c.Name()))
}
