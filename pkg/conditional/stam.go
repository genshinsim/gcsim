package conditional

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func evalStam(c *core.Core, fields []string) int64 {
	return int64(c.Player.Stam)
}
