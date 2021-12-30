package royal

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("royal spear", weapon)
	core.RegisterWeaponFunc("royalspear", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Royal(char, c, r, param)
	return "royalspear"
}
