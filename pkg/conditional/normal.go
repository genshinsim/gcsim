package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalNormalCounter(c *core.Core, fields []string) int64 {
	//.normal.bennett
	if len(fields) < 2 {
		c.Log.NewEvent("bad normal counter conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad normal counter conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad normal counter conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	return int64(char.NextNormalCounter())
}
