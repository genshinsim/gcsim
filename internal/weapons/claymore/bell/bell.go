package bell

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the bell", weapon)
	core.RegisterWeaponFunc("thebell", weapon)
}

//Taking DMG generates a shield which absorbs DMG up to 20/23/26/29/32% of Max HP.
//This shield lasts for 10s or until broken, and can only be triggered once every 45/45/45/45/45s.
//While protected by the shield, the character gains 12/15/18/21/24% increased DMG.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	hp := 0.17 + float64(r)*0.03
	icd := 0
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.09 + float64(r)*0.03

	c.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if icd > c.Frame {
			return false
		}
		icd = c.Frame + 2700 //45 seconds
		//generate a shield
		c.Shields.Add(&shield.Tmpl{
			Src:        c.Frame,
			ShieldType: core.ShieldBell,
			Name:       "Bell",
			HP:         hp * char.MaxHP(),
			Ele:        core.NoElement,
			Expires:    c.Frame + 600, //10 sec
		})
		return false
	}, fmt.Sprintf("bell-%v", char.Name()))

	char.AddMod(coretype.CharStatMod{
		Key:    "bell",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, c.Player.GetShield(core.ShieldBell) != nil
		},
	})
	return "thebell"
}
