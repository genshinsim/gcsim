package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalCD(c *core.Core, fields []string) int64 {
	//.element.t1.pyro
	if len(fields) < 3 {
		c.Log.NewEvent("bad cooldown conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	//check target is valid
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad cooldown conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad cooldown conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	var cd int

	switch fields[2] {
	case ".skill":
		cd = char.Cooldown(action.ActionSkill)
	case ".burst":
		cd = char.Cooldown(action.ActionBurst)
	default:
		c.Log.NewEvent("bad cooldown conditon: invalid action", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	//check vs the conditions
	return int64(cd)
}
