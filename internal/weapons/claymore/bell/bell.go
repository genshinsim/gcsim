package bell

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func init() {
	core.RegisterWeaponFunc("the bell", weapon)
	core.RegisterWeaponFunc("thebell", weapon)
}

//Taking DMG generates a shield which absorbs DMG up to 20/23/26/29/32% of Max HP.
//This shield lasts for 10s or until broken, and can only be triggered once every 45/45/45/45/45s.
//While protected by the shield, the character gains 12/15/18/21/24% increased DMG.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	hp := 0.17 + float64(r)*0.03
	icd := 0
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.09 + float64(r)*0.03

	c.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if icd > c.F {
			return false
		}
		icd = c.F + 2700 //45 seconds
		//generate a shield
		c.Shields.Add(&shield.Tmpl{
			Src:        c.F,
			ShieldType: core.ShieldBell,
			HP:         hp * char.MaxHP(),
			Ele:        core.NoElement,
			Expires:    c.F + 600, //10 sec
		})
		return false
	}, fmt.Sprintf("bell-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "bell",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, c.Shields.Get(core.ShieldBell) != nil
		},
	})
	return "thebell"
}
