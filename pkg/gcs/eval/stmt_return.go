package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type returnStmtEvalNode struct {
	root *ast.ReturnStmt
	node evalNode
}

func (r *returnStmtEvalNode) nextAction() (Obj, bool, error) {
	res, done, err := r.node.nextAction()
	if err != nil {
		return nil, false, err
	}
	if done {
		return &retval{
			res: res,
		}, true, nil
	}
	return res, false, nil
}

func returnStmtEval(n *ast.ReturnStmt, env *Env) evalNode {
	return &returnStmtEvalNode{
		root: n,
		node: evalFromExpr(n.Val, env),
	}
}
