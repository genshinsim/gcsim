package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalOnField(c *core.Core, fields []string) int64 {
	// .onfield.bennett
	if len(fields) < 2 {
		c.Log.NewEvent("bad onfield conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad onfield conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad onfield conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	if c.Player.Active() == char.Index {
		return 1
	}
	return 0
}
