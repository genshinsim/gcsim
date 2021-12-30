package mouunsmoon

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("mouunsmoon", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	common.Wavebreaker(char, c, r, param)
	return "mouunsmoon"
}
