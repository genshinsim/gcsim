package eval

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func binaryExprEval(n *ast.BinaryExpr, env *Env) evalNode {
	//the order of operation is going to be left, right, then root
	return &binaryExprEvalNode{
		root: n,
		stack: []evalNode{
			evalFromExpr(n.Right, env),
			evalFromExpr(n.Left, env),
		},
		res: make([]Obj, 0, 2),
	}
}

type binaryExprEvalNode struct {
	root  *ast.BinaryExpr
	res   []Obj
	stack []evalNode
}

func (b *binaryExprEvalNode) nextAction() (Obj, bool, error) {
	//eval stack while none of the results are an action
	for len(b.stack) > 0 {
		idx := len(b.stack) - 1
		res, done, err := b.stack[idx].nextAction()
		if err != nil {
			return nil, false, err
		}
		if done {
			b.stack = b.stack[:idx]
			b.res = append(b.res, res)
		}
		if res.Typ() == typAction {
			return res, false, nil //done is false b/c the binary node is not done yet
		}
	}

	//once the stack is empty, then we eval left + right
	res, err := b.evalLeftRight()
	return res, true, err
}

func (b *binaryExprEvalNode) evalLeftRight() (Obj, error) {
	if len(b.res) != 2 {
		return nil, errors.New("unexpected bool expr missing left and right")
	}
	left := b.res[0]
	right := b.res[1]

	l, ok := left.(*number)
	if !ok {
		return nil, fmt.Errorf("binary expression does not evaluate to a number, got %v ", left.Inspect())
	}
	r, ok := right.(*number)
	if !ok {
		return nil, fmt.Errorf("binary expression does not evaluate to a number, got %v ", right.Inspect())
	}
	switch b.root.Op.Typ {
	case ast.LogicAnd:
		return and(l, r), nil
	case ast.LogicOr:
		return or(l, r), nil
	case ast.ItemPlus:
		return add(l, r), nil
	case ast.ItemMinus:
		return sub(l, r), nil
	case ast.ItemAsterisk:
		return mul(l, r), nil
	case ast.ItemForwardSlash:
		return div(l, r), nil
	case ast.OpGreaterThan:
		return gt(l, r), nil
	case ast.OpGreaterThanOrEqual:
		return gte(l, r), nil
	case ast.OpEqual:
		return eq(l, r), nil
	case ast.OpNotEqual:
		return neq(l, r), nil
	case ast.OpLessThan:
		return lt(l, r), nil
	case ast.OpLessThanOrEqual:
		return lte(l, r), nil
	default:
		return &null{}, nil
	}
}
