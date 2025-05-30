package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func validateNumberParams(c *ast.CallExpr, count int) error {
	if len(c.Args) == count {
		return nil
	}
	return fmt.Errorf("invalid number of params for %v (expected %v, got %v)", c.Fun, count, len(c.Args))
}

func (e *Eval) validateArgument(c *ast.CallExpr, env *Env, index int, objType ObjTyp) (Obj, error) {
	obj, err := e.evalExpr(c.Args[index], env)
	if err != nil {
		return nil, err
	}
	if obj.Typ() != objType {
		return nil, fmt.Errorf("%v argument #%v should evaluate to %v, got %v", c.Fun, index+1, objType, obj.Typ())
	}
	return obj, nil
}

func (e *Eval) validateArguments(c *ast.CallExpr, env *Env, objTypes ...ObjTyp) ([]Obj, error) {
	err := validateNumberParams(c, len(objTypes))
	if err != nil {
		return nil, err
	}
	if len(objTypes) == 0 {
		return nil, nil
	}

	objs := make([]Obj, 0, len(objTypes))
	for i, objType := range objTypes {
		obj, err := e.validateArgument(c, env, i, objType)
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
