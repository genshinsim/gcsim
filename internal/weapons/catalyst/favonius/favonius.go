package favonius

import (
	"github.com/genshinsim/gcsim/internal/weapons/common"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("favonius codex", weapon)
	core.RegisterWeaponFunc("favoniuscodex", weapon)
	core.RegisterWeaponFunc("favcodex", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	common.Favonius(char, c, r, param)
	return "favoniuscodex"
}
