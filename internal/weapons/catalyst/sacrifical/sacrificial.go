package sacrificial

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("sacrificialfragments", weapon)
	core.RegisterWeaponFunc("sacfrags", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Sacrificial(char, c, r, param)
	return "sacrificialfragments"
}
