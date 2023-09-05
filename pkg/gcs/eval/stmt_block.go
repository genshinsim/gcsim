package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type blockStmtEvalNode struct {
	root *ast.BlockStmt
	idx  int
	env  *Env
}

func (b *blockStmtEvalNode) nextAction(env *Env) (Obj, bool, error) {
	//the first time this gets call, we should set up the block stmt env
	if b.env == nil {
		b.env = NewEnv(env)
	}
	var res Obj
	var done bool
	var err error
	for b.idx < len(b.root.List) {
		node := evalFromNode(b.root.List[b.idx])
		res, done, err = node.nextAction(b.env)
		if err != nil {
			return nil, false, err
		}
		if done {
			b.idx++
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
	return res, true, nil
}

func blockStmtEval(n *ast.BlockStmt) evalNode {
	return &blockStmtEvalNode{
		root: n,
	}
}
