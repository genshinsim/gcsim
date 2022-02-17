package queue

import (
	"errors"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalDebuff(cond core.Condition) (bool, error) {
	//.debuff.res.1.name
	if len(cond.Fields) < 4 {
		return false, errors.New("eval debuff: unexpected short field, expected at least 3")
	}
	typ := strings.TrimPrefix(cond.Fields[1], ".")
	trg := strings.TrimPrefix(cond.Fields[2], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		return false, errors.New("eval debuff: expected int for target, got " + trg)
	}

	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	active := 0
	d := strings.TrimPrefix(cond.Fields[3], ".")
	//expecting the value to be either 0 or not 0; 0 for false

	switch typ {
	case "res":
		if q.core.Combat.TargetHasResMod(d, int(tid)) {
			active = 1
		}
	case "def":
		if q.core.Combat.TargetHasDefMod(d, int(tid)) {
			active = 1
		}
	default:
		return false, nil
	}

	return compInt(cond.Op, active, val), nil
}
