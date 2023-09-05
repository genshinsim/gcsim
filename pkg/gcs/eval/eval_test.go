package eval

import (
	"errors"
	"fmt"
)

func runEvalReturnResWhenDone(e evalNode, env *Env) (Obj, error) {
	if e == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.evalNext(env)
		if err != nil {
			return nil, err
		}
		fmt.Println(val)
	}
	return val, nil
}

func runEvaluatorReturnResWhenDone(e *Evaluator) (Obj, error) {
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.base.evalNext(e.env)
		if err != nil {
			return nil, err
		}
		fmt.Println(val)
	}
	return val, nil
}
