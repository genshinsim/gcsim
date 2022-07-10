package conditional

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func evalDebuff(c *core.Core, fields []string) int64 {
	//.debuff.res.t1.name
	if len(fields) < 4 {
		c.Log.NewEvent("bad debuff conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	typ := strings.TrimPrefix(fields[1], ".")
	trg := strings.TrimPrefix(fields[2], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		c.Log.NewEvent("bad debuff conditon: invalid target", glog.LogWarnings, -1).Write("fields", fields)
	}

	d := strings.TrimPrefix(fields[3], ".")
	t := c.Combat.Target(int(tid))
	if t == nil {
		c.Log.NewEvent("bad debuff conditon: invalid target", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}
	enemy, ok := t.(*enemy.Enemy)
	if !ok {
		c.Log.NewEvent("bad debuff conditon: target not an enemy", glog.LogWarnings, -1).Write("fields", fields)
		return 0
	}

	switch typ {
	case "res":
		if enemy.ResistModIsActive(d) {
			return 1
		}
	case "def":
		if enemy.DefModIsActive(d) {
			return 1
		}
	}
	return 0
}
