package dullblade

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("dullblade", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	return "dullblade"
}
