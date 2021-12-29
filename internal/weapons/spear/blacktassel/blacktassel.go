package blacktassel

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("blacktassel", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	return "blacktassel"
}
