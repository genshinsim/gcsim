package lithic

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("lithic spear", weapon)
	core.RegisterWeaponFunc("lithicspear", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Lithic(char, c, r, param)
	return "lithicspear"
}
