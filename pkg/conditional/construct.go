package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func evalConstruct(c *core.Core, fields []string) int64 {
	if len(fields) < 3 {
		c.Log.NewEvent("bad construct conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	switch fields[1] {
	case ".duration":
		return evalConstructDuration(c, fields)
	case ".count":
		return evalConstructCount(c, fields)
	default:
		c.Log.NewEvent("bad construct conditon: invalid critera", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
}

func evalConstructDuration(c *core.Core, fields []string) int64 {
	//.construct.duration.<name>
	s := strings.TrimPrefix(fields[2], ".")
	key, ok := construct.ConstructNameToKey[s]
	if !ok {
		c.Log.NewEvent("bad construct conditon: invalid construct", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	return int64(c.Constructs.Expiry(key))
}

func evalConstructCount(c *core.Core, fields []string) int64 {
	//.construct.count.<name>
	s := strings.TrimPrefix(fields[2], ".")
	key, ok := construct.ConstructNameToKey[s]
	if !ok {
		c.Log.NewEvent("bad construct conditon: invalid construct", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	return int64(c.Constructs.CountByType(key))
}
