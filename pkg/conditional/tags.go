package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalTags(c *core.Core, fields []string) int64 {
	//.tags.char.tag
	if len(fields) < 3 {
		c.Log.NewEvent("bad tags conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	//check target is valid
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad tags conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad tags conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	tag := strings.TrimPrefix(fields[2], ".")

	return int64(char.Tag(tag))
}
