package akuoumaru

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("akuoumaru", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	common.Wavebreaker(char, c, r, param)
	return "akuoumaru"
}
