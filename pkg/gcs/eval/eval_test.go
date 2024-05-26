package eval

import (
	"errors"
)

func runEvalReturnResWhenDone(e evalNode) (Obj, error) {
	if e == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	val, _, err := e.nextAction()
	if err != nil {
		return nil, err
	}
	return val, nil
}
