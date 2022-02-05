package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalEnergy(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 2 {
		return false, errors.New("eval energy: unexpected short field, expected at least 2")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := core.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval energy: invalid char in condition")
	}
	e := char.CurrentEnergy()
	return compFloat(cond.Op, e, float64(cond.Value)), nil
}
