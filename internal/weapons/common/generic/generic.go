package generic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("generic catalyst", weapon)
	combat.RegisterWeaponFunc("generic bow", weapon)
	combat.RegisterWeaponFunc("generic claymore", weapon)
	combat.RegisterWeaponFunc("generic sword", weapon)
	combat.RegisterWeaponFunc("generic spear", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

}
