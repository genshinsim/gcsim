package blackcliff

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("blackcliffagate", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Blackcliff(char, c, r, param)
	return "blackcliffagate"
}
