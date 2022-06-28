package conditional

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func evalElement(c *core.Core, fields []string) int64 {
	//.element.t1.pyro
	if len(fields) < 3 {
		c.Log.NewEvent("bad element conditon: invalid num of fields", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	trg := strings.TrimPrefix(fields[1], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		c.Log.NewEvent("bad element conditon: invalid target", glog.LogWarnings, -1, "fields", fields)
	}
	ele := strings.TrimPrefix(fields[2], ".")
	elekey := attributes.StringToEle(ele)
	if elekey == attributes.UnknownElement {
		c.Log.NewEvent("bad element conditon: invalid element", glog.LogWarnings, -1, "fields", fields)
		return 0
	}

	t := c.Combat.Target(int(tid))
	if t == nil {
		c.Log.NewEvent("bad element conditon: invalid target", glog.LogWarnings, -1, "fields", fields)
		return 0
	}
	enemy, ok := t.(*enemy.Enemy)
	if !ok {
		c.Log.NewEvent("bad element conditon: target not an enemy", glog.LogWarnings, -1, "fields", fields)
		return 0
	}

	if enemy.AuraContains(elekey) {
		return 1
	}

	return 0
}
