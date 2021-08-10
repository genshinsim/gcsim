package bell

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/shield"
)

func init() {
	combat.RegisterWeaponFunc("the bell", weapon)
}

//Taking DMG generates a shield which absorbs DMG up to 20/23/26/29/32% of Max HP.
//This shield lasts for 10s or until broken, and can only be triggered once every 45/45/45/45/45s.
//While protected by the shield, the character gains 12/15/18/21/24% increased DMG.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	hp := 0.17 + float64(r)*0.03
	icd := 0
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.09 + float64(r)*0.03

	s.AddOnHurt(func(s core.Sim) {
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 2700 //45 seconds
		//generate a shield
		s.AddShield(&shield.Tmpl{
			Src:        s.Frame(),
			ShieldType: core.ShieldBell,
			HP:         hp * c.MaxHP(),
			Ele:        core.NoElement,
			Expires:    s.Frame() + 600, //10 sec
		})
	})

	c.AddMod(core.CharStatMod{
		Key:    "bell",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, s.GetShield(core.ShieldBell) != nil
		},
	})
}
