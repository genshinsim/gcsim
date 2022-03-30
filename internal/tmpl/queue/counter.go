package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalNormal(cond core.Condition) (bool, error) {
	//.normal.X (X = character)
	if len(cond.Fields) < 2 {
		return false, errors.New("eval normal counter: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := core.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval normal counter: invalid char in condition")
	}
	e := char.NextNormalCounter()
	return compInt(cond.Op, e, cond.Value), nil
}
