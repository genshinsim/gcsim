package blacktassel

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("blacktassel", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	return "blacktassel"
}
