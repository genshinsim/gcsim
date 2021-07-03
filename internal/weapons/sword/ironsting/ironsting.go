package ironsting

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("iron sting", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	expiry := 0
	atk := 0.045 + 0.015*float64(r)
	stacks := 0
	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.Element == def.Physical {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 60
		if expiry < s.Frame() {
			stacks = 0
		}
		stacks++
		if stacks > 2 {
			stacks = 2
		}
		expiry = s.Frame() + 360
	}, fmt.Sprintf("ironsting-%v", c.Name()))

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "ironsting",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if expiry < s.Frame() {
				stacks = 0
				return nil, false
			}
			val[def.DmgP] = atk * float64(stacks)
			return val, true
		},
	})

}
