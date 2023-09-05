package eval

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func callExprEval(n *ast.CallExpr) evalNode {
	r := &callExprEvalNode{
		root: n,
		args: make([]Obj, 0, len(n.Args)),
	}
	for _, v := range n.Args {
		r.stack = append(r.stack, evalFromExpr(v))
	}
	return r
}

type callExprEvalNode struct {
	root  *ast.CallExpr
	fn    evalNode
	args  []Obj
	stack []evalNode
}

func (c *callExprEvalNode) nextAction(env *Env) (Obj, bool, error) {
	//eval the args stack while none of the args results contains an action
	for len(c.stack) > 0 {
		idx := len(c.stack) - 1
		res, done, err := c.stack[idx].nextAction(env)
		if err != nil {
			return nil, false, err
		}
		if done {
			c.stack = c.stack[:idx]
			c.args = append(c.args, res)
		}
		if res.Typ() == typAction {
			return res, false, nil
		}
	}
	//initialize function if needed
	if c.fn == nil {
		c.fn = evalFromExpr(c.root.Fun)
	}
	//eval the expr that should return our res
	res, done, err := c.fn.nextAction(env)
	if err != nil {
		return nil, false, err
	}
	if done {
		//handle fn call only when expr is done evaluating
		return c.handleFnCall(res, env)
	}
	//the only time it's not done is if the res is an action
	if res.Typ() == typAction {
		return res, false, nil
	}
	return nil, false, fmt.Errorf("unexpected error; call expr does not evaluate to a function: %v", c.root.String())
}

func (c *callExprEvalNode) handleFnCall(fn Obj, env *Env) (Obj, bool, error) {
	switch f := fn.(type) {
	case *funcval:
		return c.handleUserFnCall(f, env)
	case *bfuncval:
		return c.handleSysFnCall(f, env)
	default:
		return nil, false, fmt.Errorf("invalid function call %v", c.root.Fun.String())
	}
}

func (c *callExprEvalNode) handleSysFnCall(fn *bfuncval, env *Env) (Obj, bool, error) {
	res, err := fn.Body(c.args, env)
	if err != nil {
		return nil, false, err
	}
	return res, true, nil
}

func (c *callExprEvalNode) handleUserFnCall(fn *funcval, env *Env) (Obj, bool, error) {
	return nil, false, errors.New("not implemented")
}
