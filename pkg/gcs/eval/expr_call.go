package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func callExprEval(n *ast.CallExpr, env *Env) evalNode {
	r := &callExprEvalNode{
		root:      n,
		args:      make([]Obj, 0, len(n.Args)),
		parentEnv: env,
	}
	for _, v := range n.Args {
		r.stack = append(r.stack, evalFromExpr(v, env))
	}
	return r
}

type callExprEvalNode struct {
	root       *ast.CallExpr
	parentEnv  *Env
	stack      []evalNode // stack is the nodes we have to eval to get the args
	args       []Obj      // track finished args
	fnCallNode evalNode   // eval into the call expr since we can have somefn()()
	fn         Obj
	fnBody     evalNode
}

func (c *callExprEvalNode) nextAction() (Obj, bool, error) {
	// eval the args stack while none of the args results contains an action
	for len(c.stack) > 0 {
		// we need to evaluate left to right
		res, done, err := c.stack[0].nextAction()
		if err != nil {
			return nil, false, err
		}
		if done {
			c.stack = c.stack[1:]
			c.args = append(c.args, res)
		}
		if res.Typ() == typAction {
			return res, false, nil
		}
	}
	//nolint:nestif // ignoring
	if c.fn == nil {
		// initialize function if needed
		if c.fnCallNode == nil {
			c.fnCallNode = evalFromExpr(c.root.Fun, c.parentEnv)
		}
		// eval the expr that should return our res
		res, done, err := c.fnCallNode.nextAction()
		if err != nil {
			return nil, false, err
		}
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; call expr does not evaluate to a function: %v", c.root.String())
		}
		// handle fn call only when expr is done evaluating
		c.fn = res
	}
	return c.handleFnCall(c.fn)
}

func (c *callExprEvalNode) handleFnCall(fn Obj) (Obj, bool, error) {
	switch f := fn.(type) {
	case *funcval:
		return c.handleUserFnCall(f)
	case *bfuncval:
		return c.handleSysFnCall(f)
	default:
		return nil, false, fmt.Errorf("invalid function call %v", c.root.Fun.String())
	}
}

func (c *callExprEvalNode) handleSysFnCall(fn *bfuncval) (Obj, bool, error) {
	res, err := fn.Body(c.args)
	if err != nil {
		return nil, false, err
	}
	return res, true, nil
}

func (c *callExprEvalNode) handleUserFnCall(fn *funcval) (Obj, bool, error) {
	if c.fnBody == nil {
		// functions are just blocks to be evaluated, along with args that are injected into the env
		for i, v := range c.args {
			fn.Env.put(fn.Args[i].Value, &v)
		}
		c.fnBody = evalFromStmt(fn.Body, fn.Env)
	}
	return c.fnBody.nextAction()
}
