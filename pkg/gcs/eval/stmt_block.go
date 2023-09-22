package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type blockStmtEvalNode struct {
	root  *ast.BlockStmt
	stack []evalNode
	env   *Env
}

func (b *blockStmtEvalNode) nextAction() (Obj, bool, error) {
	var res Obj
	var done bool
	var err error
	for len(b.stack) > 0 {
		res, done, err = b.stack[0].nextAction()
		if err != nil {
			return nil, false, err
		}
		if done {
			b.stack = b.stack[1:]
		}
		switch res.Typ() {
		case typAction:
			return res, false, nil
		case typRet:
			//if res is a return statement, then forcefully exit block regardless of
			//idx position
			return res, true, nil
		}
	}
	if res == nil {
		//this is necessary because if a block contains all actions, then block may get called again even if
		//nothing is left, resulting in a nil res
		res = &null{}
	}
	return res, true, nil
}

func blockStmtEval(n *ast.BlockStmt, env *Env) evalNode {
	b := &blockStmtEvalNode{
		root: n,
		env:  NewEnv(env),
	}
	for _, v := range n.List {
		b.stack = append(b.stack, evalFromNode(v, b.env))
	}
	return b
}
