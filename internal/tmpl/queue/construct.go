package queue

import (
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (q *Queuer) evalConstruct(cond core.Condition) (bool, error) {
	if len(cond.Fields) < 3 {
		return false, errors.New("eval construct: unexpected short field, expected at least 3")
	}
	switch cond.Fields[1] {
	case ".duration":
		return q.evalConstructDuration(cond)
	case ".count":
		return q.evalConstructCount(cond)
	default:
		return false, errors.New("eval construct: invalid field: " + cond.Fields[1])
	}
}

func (q *Queuer) evalConstructDuration(cond core.Condition) (bool, error) {
	//.construct.duration.<name>
	s := strings.TrimPrefix(cond.Fields[2], ".")
	key := core.ConstructNameToKey[s]

	//grab construct
	exp := q.core.Constructs.Expiry(key)

	return compInt(cond.Op, exp, cond.Value), nil
}

func (q *Queuer) evalConstructCount(cond core.Condition) (bool, error) {
	//.construct.count.<name>
	s := strings.TrimPrefix(cond.Fields[2], ".")
	key := core.ConstructNameToKey[s]

	//grab construct
	count := q.core.Constructs.CountByType(key)

	return compInt(cond.Op, count, cond.Value), nil
}
