package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
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
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)

	incAtk := .03 + float64(r)*0.01
	incSpd := 0.009 + float64(r)*0.003

	stacks := 0
	maxStacks := 4
	stackExpiry := 0.0
	stackDuration := 60.0 * 6.0 // frames = 6s * 60 fps

	cd := .3 * 60 // frames = 0.3s * 60fps
	icd := 0.0

	char.AddMod(core.CharStatMod{
		Key: "compoundbow",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		// Check if char is correct?
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		// Only apply on normal or charged attacks
		if (atk.Info.AttackTag != core.AttackTagNormal) && (atk.Info.AttackTag != core.AttackTagExtra) {
			return false
		}

		// Check if cd is up
		if icd > float64(c.F) {
			return false
		}

		// Assuming all stacks fall off after 6s
		if float64(c.F) > stackExpiry {
			stacks = 0
		}

		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		// trigger cd
		icd = float64(c.F) + cd
		stackExpiry = float64(c.F) + stackDuration

		m[core.ATKP] = incAtk * float64(stacks)
		m[core.AtkSpd] = incSpd * float64(stacks)

		return false
	}, fmt.Sprintf("compoundbow-%v", char.Name()))

	return "compoundbow"
}
