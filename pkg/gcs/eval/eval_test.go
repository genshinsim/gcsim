package eval

import (
	"errors"
	"fmt"
)

func runEvalReturnResWhenDone(e evalNode) (Obj, error) {
	if e == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.evalNext(nil)
		if err != nil {
			return nil, err
		}
		fmt.Println(val)
	}
	return val, nil
}
