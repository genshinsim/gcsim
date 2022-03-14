package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("compoundbow", weapon)
	core.RegisterWeaponFunc("compound-bow", weapon)
	core.RegisterWeaponFunc("compound", weapon)
}

/*
* Normal Attack and Charged Attack hits increase ATK by 4/5/6/7/8% and Normal ATK SPD by
* 1.2/1.5/1.8/2.1/2.4% for 6s. Max 4 stacks. Can only occur once every 0.3s.
 */
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)

	incAtk := .03 + float64(r)*0.01
	incSpd := 0.009 + float64(r)*0.003

	stacks := 0
	maxStacks := 4
	stackExpiry := 0
	stackDuration := 360 // frames = 6s * 60 fps

	cd := 18 // frames = 0.3s * 60fps
	icd := 0

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

		// Only apply on normal or charged attacks
		if (atk.Info.AttackTag != coretype.AttackTagNormal) && (atk.Info.AttackTag != coretype.AttackTagExtra) {
			return false
		}

		// Check if cd is up
		if icd > c.Frame {
			return false
		}

		// Reset stacks if they've expired
		if c.Frame > stackExpiry {
			stacks = 0
		}

		// Checks done, proc weapon passive
		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		// trigger cd
		icd = c.Frame + cd
		stackExpiry = c.Frame + stackDuration

		char.AddMod(coretype.CharStatMod{
			Key: "compoundbow",
			Amount: func() ([]float64, bool) {
				// Reset stacks if they've expired
				if c.Frame > stackExpiry {
					stacks = 0
				}

				m[core.ATKP] = incAtk * float64(stacks)
				m[core.AtkSpd] = incSpd * float64(stacks)
				return m, true
			},
			Expiry: stackExpiry,
		})

		return false
	}, fmt.Sprintf("compoundbow-%v", char.Name()))

	return "compoundbow"
}
