package skyrider

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyrider greatsword", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	atk := 0.05 + float64(r)*0.01
	stacks := 0
	icd := 0
	duration := 0

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
		if duration < s.Frame() {
			stacks = 0
		}

		stacks++
		if stacks > 4 {
			stacks = 4
		}
		icd = s.Frame() + 30

	}, fmt.Sprintf("skyrider-greatsword-%v", c.Name()))

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if duration > s.Frame() {
				val[def.ATKP] = atk * float64(stacks)
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})

}
