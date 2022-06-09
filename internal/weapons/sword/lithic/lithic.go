package lithic

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("lithic blade", weapon)
	core.RegisterWeaponFunc("lithicblade", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Lithic(char, c, r, param)
	return "lithicblade"
}
