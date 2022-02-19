package wine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("wine and song", weapon)
	core.RegisterWeaponFunc("wineandsong", weapon)
}

// Hitting an opponent with a Normal Attack decreases the Stamina consumption of Sprint or Alternate sprint by 14/16/18/20/22% for 5s.
// Additionally, using a Sprint or Alternate Sprint ability increases ATK by 20/25/30/35/40% for 5s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .15 + float64(r)*.05
	stam := .12 + float64(r)*.02

	c.Events.Subscribe(core.PreDash, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		char.AddMod(core.CharStatMod{
			Key:    "wineandsong",
			Expiry: c.F + 60*5,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("wineandsong-%v", char.Name()))

	stamExpiry := -1

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)

		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal {
			return false
		}

		stamExpiry = c.F + 60*5
		return false
	}, fmt.Sprintf("wineandsong-%v", char.Name()))

	c.AddStamMod(func(a core.ActionType) (float64, bool) {
		if a == core.ActionDash && stamExpiry > c.F {
			return -stam, false
		}
		return 0, false
	}, fmt.Sprintf("wineandsong-%v", char.Name()))

	return "wineandsong"
}
