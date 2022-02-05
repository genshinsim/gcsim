package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalCD(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval cd: unexpected short field, expected at least 3")
	}
	//check target is valid
	name := strings.TrimPrefix(cond.Fields[1], ".")
	key := core.CharNameToKey[name]
	char, ok := q.core.CharByName(key)
	if !ok {
		return false, errors.New("eval cd: invalid char in condition")
	}
	var cd int
	switch cond.Fields[2] {
	case ".skill":
		cd = char.Cooldown(core.ActionSkill)
	case ".burst":
		cd = char.Cooldown(core.ActionBurst)
	default:
		return false, nil
	}
	//check vs the conditions
	return compInt(cond.Op, cd, cond.Value), nil
}
