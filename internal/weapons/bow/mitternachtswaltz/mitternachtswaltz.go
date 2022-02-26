package mitternachtswaltz

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
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
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)

	buffAmount := .15 + .05*float64(r) // same amount in either context
	buffExpiry := 300                  // 5s

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		// Attack belongs to the equipped character
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		// Active character has weapon equipped
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		// only apply elemental skill buff on normal attacks
		if atk.Info.AttackTag == core.AttackTagNormal {
			char.AddPreDamageMod(core.PreDamageMod{
				Key: "mitternachtswaltz-ele",
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					if (atk.Info.AttackTag == core.AttackTagElementalArt) || (atk.Info.AttackTag == core.AttackTagElementalArtHold) {
						m[core.DmgP] = buffAmount
						return m, true
					}

					return nil, false
				},
				Expiry: buffExpiry,
			})
		}

		// only apply normal attack buff on elemental skill
		if (atk.Info.AttackTag == core.AttackTagElementalArt) || (atk.Info.AttackTag == core.AttackTagElementalArtHold) {
			char.AddPreDamageMod(core.PreDamageMod{
				Key: "mitternachtswaltz-na",
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					if atk.Info.AttackTag == core.AttackTagNormal {
						m[core.DmgP] = buffAmount
						return m, true
					}

					return nil, false
				},
				Expiry: buffExpiry,
			})
		}

		return false
	}, fmt.Sprintf("mitternachtswaltz-%v", char.Name()))

	return "mitternachtswaltz"
}
