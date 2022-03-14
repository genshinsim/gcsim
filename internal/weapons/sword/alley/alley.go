package alley

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the alley flash", weapon)
	core.RegisterWeaponFunc("thealleyflash", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	lockout := -1

	c.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		lockout = c.Frame + 300
		return false
	}, fmt.Sprintf("alley-flash-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.09 + 0.03*float64(r)
	char.AddMod(coretype.CharStatMod{
		Key: "royal",
		Amount: func() ([]float64, bool) {
			return m, lockout < c.Frame
		},
		Expiry: -1,
	})
	return "thealleyflash"
}
