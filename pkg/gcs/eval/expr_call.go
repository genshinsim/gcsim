package eval

import (
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
	root       *ast.CallExpr
	fnCallNode evalNode
	fn         Obj
	fnBody     evalNode
	fnEnv      *Env //tracking this here to avoid rebuilding this every time
	args       []Obj
	stack      []evalNode
}

func (c *callExprEvalNode) nextAction(env *Env) (Obj, bool, error) {
	//eval the args stack while none of the args results contains an action
	for len(c.stack) > 0 {
		//we need to evaluate left to right
		res, done, err := c.stack[0].nextAction(env)
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
	if c.fn == nil {
		//initialize function if needed
		if c.fnCallNode == nil {
			c.fnCallNode = evalFromExpr(c.root.Fun)
		}
		//eval the expr that should return our res
		res, done, err := c.fnCallNode.nextAction(env)
		if err != nil {
			return nil, false, err
		}
		if !done {
			//the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; call expr does not evaluate to a function: %v", c.root.String())
		}
		//handle fn call only when expr is done evaluating
		c.fn = res
	}
	return c.handleFnCall(c.fn, env)
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
	//functions are just blocks to be evaluated, along with args that are injected into the env
	if c.fnEnv == nil {
		c.fnEnv = NewEnv(fn.Env)
		for i, v := range c.args {
			c.fnEnv.put(fn.Args[i].Value, &v)
		}
	}
	if c.fnBody == nil {
		c.fnBody = evalFromStmt(fn.Body)
	}
	return c.fnBody.nextAction(c.fnEnv)
}
