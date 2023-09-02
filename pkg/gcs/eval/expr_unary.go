package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type unaryExprEvalNode struct {
	root *ast.UnaryExpr
	node evalNode
}

func (u *unaryExprEvalNode) evalNext(env *Env) (Obj, bool, error) {
	if u.node == nil {
		u.node = evalFromExpr(u.root.Right)
	}
	res, done, err := u.node.evalNext(env)
	if err != nil {
		return nil, false, err
	}
	//check not done first because once done we'll need to apply the unary
	//operation; cleaner this way
	if !done {
		return res, false, err
	}
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

func unaryExprEval(n *ast.UnaryExpr) evalNode {
	return &unaryExprEvalNode{
		root: n,
	}
}
