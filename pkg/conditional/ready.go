package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalReady(c *core.Core, fields []string) int64 {
	if len(fields) < 3 {
		c.Log.NewEvent("bad ready conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	//check target is valid
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad ready conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad ready conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}

	abil := strings.TrimPrefix(fields[2], ".")
	ak := action.StringToAction(abil)
	if ak == action.InvalidAction {
		c.Log.NewEvent("bad ready conditon: invalid action", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}

	//TODO: nil map may cause problems here??
	if char.ActionReady(ak, nil) {
		return 1
	}

	return 0
}
