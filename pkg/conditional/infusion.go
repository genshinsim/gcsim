package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalInfusion(c *core.Core, fields []string) int64 {
	//.infusion.bennett.key
	if len(fields) < 3 {
		c.Log.NewEvent("bad infusion conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad infusion conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad infusion conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	inf := strings.TrimPrefix(fields[2], ".")
	if c.Player.WeaponInfuseIsActive(char.Index, inf) {
		return 1
	}
	return 0
}
