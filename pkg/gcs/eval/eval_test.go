package eval

import (
	"errors"
)

func runEvalReturnResWhenDone(e evalNode, env *Env) (Obj, error) {
	if e == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	val, _, err := e.nextAction(env)
	if err != nil {
		return nil, err
	}
	return val, nil
}
