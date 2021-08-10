package generic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("generic catalyst", weapon)
	combat.RegisterWeaponFunc("generic bow", weapon)
	combat.RegisterWeaponFunc("generic claymore", weapon)
	combat.RegisterWeaponFunc("generic sword", weapon)
	combat.RegisterWeaponFunc("generic spear", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

}
