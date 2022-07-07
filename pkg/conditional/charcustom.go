package conditional

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func evalCharCustom(c *core.Core, key keys.Char, fields []string) int64 {
	if len(fields) < 2 {
		c.Log.NewEvent("bad char custom conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		return 0
	}
	return char.Condition(fields[1])
}
