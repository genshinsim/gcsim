package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func evalStatus(c *core.Core, fields []string) int64 {
	//.energy.char
	if len(fields) < 2 {
		c.Log.NewEvent("bad status conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	//check target is valid
	status := strings.TrimPrefix(fields[1], ".")

	return int64(c.Status.Duration(status))
}
