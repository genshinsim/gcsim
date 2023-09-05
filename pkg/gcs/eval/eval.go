package eval

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type evalNode interface {
	//execute the node; node should either return next action, or continue execution until the node is
	//done. done should only be false if Obj is an action; otherwise must be true
	nextAction(*Env) (Obj, bool, error)
}

type Evaluator struct {
	Core *core.Core
	base evalNode
	env  *Env
}

func NewEvaluator(root ast.Node) (*Evaluator, error) {
	e := &Evaluator{
		env: NewEnv(nil),
	}
	e.base = evalFromNode(root)
	if e.base == nil {
		return nil, errors.New("invalid root node; no executor found")
	}
	return e, nil
}

func (e *Evaluator) NextAction() (*action.ActionEval, error) {
	//continue eval until we hit an action
	for {
		//base case: no more action
		if e.base == nil {
			return nil, nil
		}
		res, done, err := e.base.nextAction(e.env)
		if err != nil {
			return nil, err
		}
		if done {
			e.base = nil
		}
		//we're done
		if v, ok := res.(*actionval); ok {
			return v.toActionEval(), nil
		}
	}
}

func evalFromNode(n ast.Node) evalNode {
	switch v := n.(type) {
	case ast.Expr:
		return evalFromExpr(v)
	case ast.Stmt:
		return evalFromStmt(v)
	default:
		return nil
	}
}
