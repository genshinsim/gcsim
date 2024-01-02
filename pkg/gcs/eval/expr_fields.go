package eval

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func fieldsExprEval(n *ast.Field, env *Env) evalNode {
	r := &fieldsExprEvalNode{
		parentEnv: env,
		root:      n,
	}
	for _, v := range n.Value {
		r.args = append(r.args, &strval{str: v})
	}
	return r
}

type fieldsExprEvalNode struct {
	root      *ast.Field
	parentEnv *Env
	args      []Obj
	fnBody    evalNode // this is to allow override of evaluate_conditional
}

func (c *fieldsExprEvalNode) nextAction() (Obj, bool, error) {
	fn, ok := c.parentEnv.v("evaluate_conditional")
	if !ok {
		return nil, false, errors.New("unexpected system error, sysfunc evaluate_conditional not available")
	}
	return c.handleFnCall(*fn)
}

func (c *fieldsExprEvalNode) handleFnCall(fn Obj) (Obj, bool, error) {
	switch f := fn.(type) {
	case *funcval:
		return c.handleUserFnCall(f)
	case *bfuncval:
		return c.handleSysFnCall(f)
	default:
		return nil, false, fmt.Errorf("unexpected system error, evaluate_coditional should be a function but got %v", fn.Typ())
	}
}

func (c *fieldsExprEvalNode) handleSysFnCall(fn *bfuncval) (Obj, bool, error) {
	res, err := fn.Body(c.args)
	if err != nil {
		return nil, false, err
	}
	return res, true, nil
}

func (c *fieldsExprEvalNode) handleUserFnCall(fn *funcval) (Obj, bool, error) {
	if c.fnBody == nil {
		// functions are just blocks to be evaluated, along with args that are injected into the env
		for i, v := range c.args {
			fn.Env.put(fn.Args[i].Value, &v)
		}
		c.fnBody = evalFromStmt(fn.Body, fn.Env)
	}
	return c.fnBody.nextAction()
}
