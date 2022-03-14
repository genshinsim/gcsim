package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (q *Queuer) evalTags(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval tags: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := core.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval tags: invalid char in condition")
	}
	tag := strings.TrimPrefix(cond.Fields[2], ".")
	v := char.Tag(tag)
	q.coretype.Log.NewEvent("evaluating tags", coretype.LogQueueEvent, char.Index(), "targ", tag, "val", v)
	return compInt(cond.Op, v, cond.Value), nil
}
