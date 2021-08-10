package solar

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("solar pearl", weapon)
}

//Normal Attack hits increase Elemental Skill and Elemental Burst DMG by 20/25/30/35/40% for 6s.
//Likewise, Elemental Skill or Elmental Burst hits increase Normal Attack DMG by 20/25/30/35/40% for 6s.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	skill := 0
	attack := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag == core.AttackTagElementalArt || ds.AttackTag == core.AttackTagElementalBurst {
			skill = s.Frame() + 300
			return
		}
		if ds.AttackTag == core.AttackTagNormal {
			skill = s.Frame() + 300
		}
	}, fmt.Sprintf("solar-%v", c.Name()))

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15 + float64(r)*0.05
	c.AddMod(core.CharStatMod{
		Key:    "solar",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalArt || a == core.AttackTagElementalBurst {
				return val, attack > s.Frame()
			}
			if a == core.AttackTagNormal {
				return val, skill > s.Frame()
			}
			return nil, false
		},
	})
}
