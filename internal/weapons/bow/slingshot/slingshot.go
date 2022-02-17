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
	decrDmg := 0.01
	passiveThreshold := 18 // 0.3s

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

		// Only apply on normal or charged attacks
		if (atk.Info.AttackTag != core.AttackTagNormal) && (atk.Info.AttackTag != core.AttackTagExtra) {
			return false
		}

		dmgP := decrDmg
		arrowHangTime := c.F - atk.SourceFrame
		if arrowHangTime > passiveThreshold {
			dmgP = incrDmg
		}

		char.AddMod(core.CharStatMod{
			Key: "slingshot",
			Amount: func() ([]float64, bool) {
				m[core.DmgP] = dmgP
				return m, true
			},
			Expiry: -1,
		})

		return false
	}, fmt.Sprintf("slingshot-%v", char.Name()))

	return "slingshot"
}
