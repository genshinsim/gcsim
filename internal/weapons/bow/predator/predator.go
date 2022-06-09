package predator

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("predator", weapon)
}

/*
 * Dealing Cryo DMG to opponents increases this character's Normal and Charged Attack DMG by 10% for 6s.
 * This effect can have a maximum of 2 stacks.
 * Additionally, when Aloy equips Predator, ATK is increased by 66.
 * (Refines do not change the weapon effect)
 */
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	mATK, mDMG := make([]float64, core.EndStatType), make([]float64, core.EndStatType)

	if char.Key() == core.Aloy {
		char.AddMod(core.CharStatMod{
			Key:    "predator",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				mATK[core.ATK] = 66
				return mATK, true
			},
		})
	}

	buffDmgP := .10

	stacks := 0
	stackExpiry := 0
	maxStacks := 2
	stackDuration := 360 // 6s

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

		// Only apply when damage element is cryo
		if atk.Info.Element != core.Cryo {
			return false
		}

		// Reset stacks if they've expired
		if c.F > stackExpiry {
			stacks = 0
		}

		// Checks done, proc weapon passive
		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		stackExpiry = c.F + stackDuration

		char.AddPreDamageMod(core.PreDamageMod{
			Key: "predator",
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				// Reset stacks if they've expired
				if c.F > stackExpiry {
					stacks = 0
				}

				// Only apply to normal or charged attacks
				if (atk.Info.AttackTag == core.AttackTagNormal) || (atk.Info.AttackTag == core.AttackTagExtra) {
					mDMG[core.DmgP] = buffDmgP * float64(stacks)
					return mDMG, true
				}

				return nil, false
			},
			Expiry: stackExpiry,
		})

		return false
	}, fmt.Sprintf("predator-%v", char.Name()))

	return "predator"
}
