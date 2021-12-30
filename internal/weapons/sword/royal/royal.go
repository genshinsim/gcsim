package royal

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("royal longsword", weapon)
	core.RegisterWeaponFunc("royallongsword", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Royal(char, c, r, param)
	return "royallongsword"
}
