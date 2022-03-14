package primordial

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("primordial jade winged-spear", weapon)
	core.RegisterWeaponFunc("primordialjadewingedspear", weapon)
}

//For every character in the party who hails from Liyue, the character who equips this
//weapon gains 6/7/8/9//10% ATK increase and 2/3/4/5/6% CRIT Rate increase.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	last := 0
	stacks := 0
	active := 0

	m := make([]float64, core.EndStatType)

	char.AddMod(coretype.CharStatMod{
		Key: "primordial",
		Amount: func() ([]float64, bool) {
			return m, active > c.Frame
		},
		Expiry: -1,
	})

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		//check if char is correct?
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.ActiveChar != char.Index() {
			return false
		}
		//check if cd is up
		if c.Frame-last < 18 && last != 0 {
			return false
		}
		//check if expired; reset stacks if so
		if active < c.Frame {
			stacks = 0
		}

		stacks++
		active = c.Frame + 360

		if stacks >= 7 {
			stacks = 7
			m[core.DmgP] = 0.09 + float64(r)*0.03
		}
		m[core.ATKP] = (float64(r)*0.007 + 0.025) * float64(stacks)

		//trigger cd
		last = c.Frame
		return false
	}, fmt.Sprintf("primordial-%v", char.Name()))
	return "primordialjadewingedspear"
}
