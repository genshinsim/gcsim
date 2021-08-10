package primordial

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("primordial jade winged-spear", weapon)
}

//For every character in the party who hails from Liyue, the character who equips this
//weapon gains 6/7/8/9//10% ATK increase and 2/3/4/5/6% CRIT Rate increase.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	last := 0
	stacks := 0
	active := 0

	m := make([]float64, core.EndStatType)

	c.AddMod(core.CharStatMod{
		Key: "primordial",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, active > s.Frame()
		},
		Expiry: -1,
	})

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		//check if char is correct?
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		//check if cd is up
		if s.Frame()-last < 18 && last != 0 {
			return
		}
		//check if expired; reset stacks if so
		if active < s.Frame() {
			stacks = 0
		}

		stacks++
		active = s.Frame() + 360

		if stacks > 7 {
			stacks = 7
			m[core.DmgP] = 0.09 + float64(r)*0.03
		}
		m[core.ATK] = (float64(r)*0.007 + 0.025) * float64(stacks)

		//trigger cd
		last = s.Frame()
	}, fmt.Sprintf("primordial-%v", c.Name()))

}
