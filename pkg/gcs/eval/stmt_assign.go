package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type assignStmtEvalNode struct {
	*ast.AssignStmt
	env  *Env
	node evalNode
	res  Obj
}

func (l *assignStmtEvalNode) nextAction() (Obj, bool, error) {
	if l.res == nil {
		res, done, err := l.node.nextAction()
		if err != nil {
			return nil, false, err
		}
		if !done {
			// the only time it's not done is if the res is an action
			if res.Typ() == typAction {
				return res, false, nil
			}
			return nil, false, fmt.Errorf("unexpected error; assign stmt stopped at non action: %v", l.AssignStmt.String())
		}
		l.res = res
	}
	ok := l.env.assign(l.Ident.Val, &l.res)
	if !ok {
		return nil, false, fmt.Errorf("variable %v does not exist; cannot assign", l.Ident.Val)
	}
	return l.res, true, nil
}

func assignStmtEval(n *ast.AssignStmt, env *Env) evalNode {
	return &assignStmtEvalNode{
		AssignStmt: n,
		env:        env,
		node:       evalFromExpr(n.Val, env),
	}
}
