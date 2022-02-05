package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalStatus(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval status: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	status := q.core.Status.Duration(name)
	// q.core.Log.Debugw("queue status check", "frame", q.core.F, "event", LogQueueEvent, "status", name, "val", status, "expected", c.Value, "op", c.Op)
	return compInt(cond.Op, status, cond.Value), nil

}
