package eval

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type evalNode interface {
	// execute the node; node should either return next action, or continue execution until the node is
	// done. done should only be false if Obj is an action; otherwise must be true
	nextAction() (Obj, bool, error)
}

type Evaluator struct {
	Core *core.Core
	base evalNode
	env  *Env
	err  error
}

func NewEvaluator(root ast.Node, core *core.Core) (*Evaluator, error) {
	e := &Evaluator{
		Core: core,
		env:  NewEnv(nil),
	}
	e.initSysFuncs()
	e.base = evalFromNode(root, e.env)
	if e.base == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	return e, nil
}

func (e *Evaluator) NextAction() (*action.Eval, error) {
	// base case: no more action
	if e.base == nil {
		return nil, e.err
	}
	res, done, err := e.base.nextAction()
	if err != nil {
		e.err = err
		e.base = nil
		return nil, err
	}
	if done {
		e.base = nil
	}
	if v, ok := res.(*actionval); ok {
		return v.toActionEval(), nil
	}
	return nil, nil
}

func (e *Evaluator) Continue()   {}
func (e *Evaluator) Exit() error { return nil }
func (e *Evaluator) Err() error  { return e.err }
func (e *Evaluator) Start()      {}

func evalFromNode(n ast.Node, env *Env) evalNode {
	switch v := n.(type) {
	case ast.Expr:
		return evalFromExpr(v, env)
	case ast.Stmt:
		return evalFromStmt(v, env)
	default:
		return nil
	}
}
