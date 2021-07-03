package solar

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("solar pearl", weapon)
}

//Normal Attack hits increase Elemental Skill and Elemental Burst DMG by 20/25/30/35/40% for 6s.
//Likewise, Elemental Skill or Elmental Burst hits increase Normal Attack DMG by 20/25/30/35/40% for 6s.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	skill := 0
	attack := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag == def.AttackTagElementalArt || ds.AttackTag == def.AttackTagElementalBurst {
			skill = s.Frame() + 300
			return
		}
		if ds.AttackTag == def.AttackTagNormal {
			skill = s.Frame() + 300
		}
	}, fmt.Sprintf("solar-%v", c.Name()))

	val := make([]float64, def.EndStatType)
	val[def.DmgP] = 0.15 + float64(r)*0.05
	c.AddMod(def.CharStatMod{
		Key:    "solar",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a == def.AttackTagElementalArt || a == def.AttackTagElementalBurst {
				return val, attack > s.Frame()
			}
			if a == def.AttackTagNormal {
				return val, skill > s.Frame()
			}
			return nil, false
		},
	})
}
