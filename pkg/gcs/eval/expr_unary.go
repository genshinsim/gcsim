package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type unaryExprEvalNode struct {
	root *ast.UnaryExpr
	node evalNode
}

func (u *unaryExprEvalNode) nextAction() (Obj, bool, error) {
	res, done, err := u.node.nextAction()
	if err != nil {
		return nil, false, err
	}
	if done {
		return u.handleUnaryOperation(res)
	}
	//the only time it's not done is if the res is an action
	if res.Typ() == typAction {
		return res, false, nil
	}
	return nil, false, fmt.Errorf("unexpected error; unary expr does not evaluate to a value: %v", u.root.Right.String())
}

func (u *unaryExprEvalNode) handleUnaryOperation(res Obj) (Obj, bool, error) {
	r, ok := res.(*number)
	if !ok {
		return nil, false, fmt.Errorf("unary expression does not evaluate to a number, got %v ", res.Inspect())
	}
	switch u.root.Op.Typ {
	case ast.LogicNot:
		return eq(&number{}, r), true, nil
	case ast.ItemMinus:
		return sub(&number{}, r), true, nil
	default:
		return nil, false, fmt.Errorf("unrecognized unary operator %v", u.root.Op)
	}
}

func unaryExprEval(n *ast.UnaryExpr, env *Env) evalNode {
	return &unaryExprEvalNode{
		root: n,
		node: evalFromExpr(n.Right, env),
	}
}
