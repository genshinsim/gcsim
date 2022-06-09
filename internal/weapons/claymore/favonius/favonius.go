package favonius

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("favonius greatsword", weapon)
	core.RegisterWeaponFunc("favoniusgreatsword", weapon)
	core.RegisterWeaponFunc("favgreatsword", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Favonius(char, c, r, param)
	return "favoniusgreatsword"
}
