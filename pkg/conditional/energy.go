package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalEnergy(c *core.Core, fields []string) int64 {
	//.energy.char
	if len(fields) < 2 {
		c.Log.NewEvent("bad energy conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	//check target is valid
	name := strings.TrimPrefix(fields[1], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad energy conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		c.Log.NewEvent("bad energy conditon: invalid character", glog.LogWarnings, -1, "fields", fields)
		return 0
	}

	//this will floor it
	return int64(char.Energy)
}
