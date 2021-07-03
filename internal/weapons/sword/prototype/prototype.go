package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("prototype rancour", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	expiry := 0
	per := 0.03 + 0.01*float64(r)
	stacks := 0
	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 18
		if expiry < s.Frame() {
			stacks = 0
		}
		stacks++
		if stacks > 4 {
			stacks = 4
		}
		expiry = s.Frame() + 360
	}, fmt.Sprintf("prototype-rancour-%v", c.Name()))

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if expiry < s.Frame() {
				stacks = 0
				return nil, false
			}
			val[def.ATKP] = per * float64(stacks)
			val[def.DEFP] = per * float64(stacks)
			return val, true
		},
	})

}
