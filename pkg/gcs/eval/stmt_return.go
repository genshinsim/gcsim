package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type returnStmtEvalNode struct {
	root *ast.ReturnStmt
	node evalNode
}

func (r *returnStmtEvalNode) evalNext(env *Env) (Obj, bool, error) {
	if r.node == nil {
		r.node = evalFromExpr(r.root.Val)
	}

	res, done, err := r.node.evalNext(env)
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

func returnStmtEval(n *ast.ReturnStmt) evalNode {
	return &returnStmtEvalNode{root: n}
}
