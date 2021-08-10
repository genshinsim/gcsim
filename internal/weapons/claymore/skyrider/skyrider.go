package skyrider

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyrider greatsword", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	atk := 0.05 + float64(r)*0.01
	stacks := 0
	icd := 0
	duration := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
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

	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if duration > s.Frame() {
				val[core.ATKP] = atk * float64(stacks)
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})

}
