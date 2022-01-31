package solar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("solar pearl", weapon)
	core.RegisterWeaponFunc("solarpearl", weapon)
}

//Normal Attack hits increase Elemental Skill and Elemental Burst DMG by 20/25/30/35/40% for 6s.
//Likewise, Elemental Skill or Elmental Burst hits increase Normal Attack DMG by 20/25/30/35/40% for 6s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15 + float64(r)*0.05

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagElementalArt, core.AttackTagElementalArtHold, core.AttackTagElementalBurst:
			char.AddPreDamageMod(core.PreDamageMod{
				Key:    "solar-na-buff",
				Expiry: c.F + 6*60,
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					switch atk.Info.AttackTag {
					case core.AttackTagNormal:
						return val, true
					}
					return nil, false
				},
			})
		case core.AttackTagNormal:
			char.AddPreDamageMod(core.PreDamageMod{
				Key:    "solar-skill-burst-buff",
				Expiry: c.F + 6*60,
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					switch atk.Info.AttackTag {
					case core.AttackTagElementalArt, core.AttackTagElementalArtHold, core.AttackTagElementalBurst:
						return val, true
					}
					return nil, false
				},
			})
		default:
			return false
		}
		return false
	}, fmt.Sprintf("solar-%v", char.Name()))

	return "solarpearl"
}
