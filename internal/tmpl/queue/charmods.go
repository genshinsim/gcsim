package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalCharMods(cond core.Condition) (bool, error) {
	//.mods.bennett.buff==1
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
	val := cond.Value
	if val > 0 {
		val = 1
	} else {
		val = 0
	}
	q.core.Log.Debugw("evaluating mods", "frame", q.core.F, "event", core.LogQueueEvent, "char", char.CharIndex(), "mod", tag)
	return char.ModIsActive(tag) == (val == 1), nil
}
