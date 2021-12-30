package waster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("wastergreatsword", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	return "wastergreatsword"
}
