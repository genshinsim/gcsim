package bell

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/shield"
)

func init() {
	combat.RegisterWeaponFunc("the bell", weapon)
}

//Taking DMG generates a shield which absorbs DMG up to 20/23/26/29/32% of Max HP.
//This shield lasts for 10s or until broken, and can only be triggered once every 45/45/45/45/45s.
//While protected by the shield, the character gains 12/15/18/21/24% increased DMG.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	hp := 0.17 + float64(r)*0.03
	icd := 0
	val := make([]float64, def.EndStatType)
	val[def.DmgP] = 0.09 + float64(r)*0.03

	s.AddOnHurt(func(s def.Sim) {
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 2700 //45 seconds
		//generate a shield
		s.AddShield(&shield.Tmpl{
			Src:        s.Frame(),
			ShieldType: def.ShieldBell,
			HP:         hp * c.MaxHP(),
			Ele:        def.NoElement,
			Expires:    s.Frame() + 600, //10 sec
		})
	})

	c.AddMod(def.CharStatMod{
		Key:    "bell",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, s.GetShield(def.ShieldBell) != nil
		},
	})
}
