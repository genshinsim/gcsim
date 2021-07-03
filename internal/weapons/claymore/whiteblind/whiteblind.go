package whiteblind

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("whiteblind", weapon)
}

//On hit, Normal or Charged Attacks increase ATK and DEF by 6/7.5/9/10.5/12% for 6s.
//Max 4 stacks (24/30/36/42/48% total). This effect can only occur once every 0.5s.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
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
		duration = s.Frame() + 360
		icd = s.Frame() + 30
		stacks++
		if stacks > 4 {
			stacks = 4
		}
	}, fmt.Sprintf("whiteblind-%v", c.Name()))

	amt := 0.045 + float64(r)*0.015

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "whiteblind",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if duration < s.Frame() {
				stacks = 0
				return nil, false
			}
			val[def.ATKP] = amt * float64(stacks)
			val[def.DEFP] = amt * float64(stacks)
			return val, true
		},
	})
}
