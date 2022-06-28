package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalCharMods(c *core.Core, fields []string) int64 {
	//.mods.bennett.buff==1
	if len(fields) < 3 {
		c.Log.NewEvent("bad char mod conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad char mod conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad char mod conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	tag := strings.TrimPrefix(fields[2], ".")
	//TODO: be nice if we can check attackmods somehow but those are conditional
	//on attacks/targets and we cant really supply a fake attack or fake target here
	if char.StatModIsActive(tag) {
		return 1
	}
	return 0
}
