package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalInfusion(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval infusion: unexpected short field, expected at least 3")
	}
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := core.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval infusion: invalid char in condition")
	}

	active := 0
	inf := strings.TrimPrefix(cond.Fields[2], ".")
	if char.WeaponInfuseIsActive(inf) {
		active = 1
	}

	return compInt(cond.Op, active, cond.Value), nil
}
