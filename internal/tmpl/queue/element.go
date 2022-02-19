package queue

import (
	"errors"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalElement(cond core.Condition) (bool, error) {
	//.element.1.pyro
	if len(cond.Fields) < 3 {
		return false, errors.New("eval element: unexpected short field, expected at least 2")
	}
	trg := strings.TrimPrefix(cond.Fields[1], ".t")
	//trg should be an int
	tid, err := strconv.ParseInt(trg, 10, 64)
	if err != nil {
		//invalid target
		return false, errors.New("eval element: expected int for target, got " + trg)
	}

	ele := strings.TrimPrefix(cond.Fields[2], ".")
	//expecting the value to be either 0 or not 0; 0 for false
	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	active := 0
	e := core.StringToEle(ele)
	if e == core.UnknownElement {
		return false, nil
	}

	if q.core.Combat.TargetHasElement(e, int(tid)) {
		active = 1
	}
	return compInt(cond.Op, active, val), nil
}
