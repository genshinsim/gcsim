package memory

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("memory of dust", weapon)
	core.RegisterWeaponFunc("memoryofdust", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.GoldenMajesty(char, c, r, param)
	return "memoryofdust"
}
