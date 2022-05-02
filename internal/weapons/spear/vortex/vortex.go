package vortex

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("vortex vanquisher", weapon)
	core.RegisterWeaponFunc("vortexvanquisher", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.GoldenMajesty(char, c, r, param)
	return "vortexvanquisher"
}
