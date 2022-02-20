package generic

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("slingshot", weapon)
}

/*
* If a Normal or Charged Attack hits a target within 0.3s of being fired, increases DMG by 36/42/48/54/60%.
* Otherwise, decreases DMG by 10%.
 */
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)

	incrDmg := .3 + float64(r)*0.06
	decrDmg := -0.10
	passiveThresholdF := 18 // 0.3s
	travel := 0

	char.AddPreDamageMod(core.PreDamageMod{
		Key: "slingshot",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if (atk.Info.AttackTag != core.AttackTagNormal) && (atk.Info.AttackTag != core.AttackTagExtra) {
				return nil, false
			}

			m[core.DmgP] = incrDmg
			if travel > passiveThresholdF {
				m[core.DmgP] = decrDmg
			}

			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		// Attack belongs to the equipped character
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		// Active character has weapon equipped
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		// Only apply on normal or charged attacks
		if (atk.Info.AttackTag != core.AttackTagNormal) && (atk.Info.AttackTag != core.AttackTagExtra) {
			return false
		}

		// Calc travel for the mod
		travel = c.F - atk.Snapshot.SourceFrame

		return false
	}, fmt.Sprintf("slingshot-%v", char.Name()))

	return "slingshot"
}
