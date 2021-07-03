package memory

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("memory of dust", weapon)
}

//Increases Shield Strength by 20/25/30/35/40%. Scoring hits on opponents increases ATK by 4/5/6/7/8% for 8s. Max 5 stacks.
//Can only occur once every 0.3s. While protected by a shield, this ATK increase effect is increased by 100%.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	shd := .15 + float64(r)*.05
	s.AddShieldBonus(func() float64 {
		return shd
	})

	stacks := 0
	icd := 0
	duration := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}
		if duration < s.Frame() {
			stacks = 0
		}
		stacks++
		if stacks > 5 {
			stacks = 0
		}

	}, fmt.Sprintf("memory-dust-%v", c.Name()))

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "memory",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if duration > s.Frame() {
				val[def.ATKP] = 0.04 * float64(stacks)
				if s.IsShielded() {
					val[def.ATKP] *= 2
				}
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})

}
