package mitternachtswaltz

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("mitternachtswaltz", weapon)
	core.RegisterWeaponFunc("mitternachts", weapon)
	core.RegisterWeaponFunc("mitternacht", weapon)
}

/*
 * Normal Attack hits on opponents increase Elemental Skill DMG by 20/25/30/35/40% for 5s.
 * Elemental Skill hits on opponents increase Normal Attack DMG by 20/25/30/35/40% for 5s.
 */
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)

	buffAmount := .15 + .05*float64(r) // same amount in either context
	buffExpiry := 300                  // 5s
	buffIcd := 0                       // Add a 1-frame ICD to prevent buffs from being applied too quickly for the sim

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)

		// Attack belongs to the equipped character
		if atk.Info.ActorIndex != char.Index() {
			return false
		}

		// Active character has weapon equipped
		if c.ActiveChar != char.Index() {
			return false
		}

		// Add 1-frame ICD to prevent too many buffs from being applied the sim simultaneously
		if c.Frame <= buffIcd {
			return false
		}

		buffIcd = c.Frame + 1

		// only apply elemental skill buff on normal attacks
		if atk.Info.AttackTag == coretype.AttackTagNormal {
			char.AddPreDamageMod(coretype.PreDamageMod{
				Key: "mitternachtswaltz-ele",
				Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
					if (atk.Info.AttackTag == core.AttackTagElementalArt) || (atk.Info.AttackTag == core.AttackTagElementalArtHold) {
						m[core.DmgP] = buffAmount
						return m, true
					}

					return nil, false
				},
				Expiry: c.Frame + buffExpiry,
			})
		}

		// only apply normal attack buff on elemental skill
		if (atk.Info.AttackTag == core.AttackTagElementalArt) || (atk.Info.AttackTag == core.AttackTagElementalArtHold) {
			char.AddPreDamageMod(coretype.PreDamageMod{
				Key: "mitternachtswaltz-na",
				Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
					if atk.Info.AttackTag == coretype.AttackTagNormal {
						m[core.DmgP] = buffAmount
						return m, true
					}

					return nil, false
				},
				Expiry: c.Frame + buffExpiry,
			})
		}

		return false
	}, fmt.Sprintf("mitternachtswaltz-%v", char.Name()))

	return "mitternachtswaltz"
}
