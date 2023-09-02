package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

type blockStmtEvalNode struct {
	root        *ast.BlockStmt
	currentNode evalNode
	idx         int
	env         *Env
}

func (b *blockStmtEvalNode) evalNext(env *Env) (Obj, bool, error) {
	//the first time this gets call, we should set up the block stmt env
	if b.env == nil {
		b.env = NewEnv(env)
	}
	//if currentNode is nil, that means this idx has not been visited yet
	if b.currentNode == nil {
		b.currentNode = evalFromNode(b.root.List[b.idx])
	}
	res, done, err := b.currentNode.evalNext(b.env)
	if err != nil {
		return nil, false, err
	}
	if done {
		b.currentNode = nil
		b.idx++
	}
	//if res is a return statement, then forcefully exit block regardless of
	//idx position
	if res.Typ() == typRet {
		return res, true, nil
	}
	//otherwise block is done if idx == lenght
	if b.idx == len(b.root.List) {
		return res, true, nil
	}
	return res, false, nil
}

func blockStmtEval(n *ast.BlockStmt) evalNode {
	return &blockStmtEvalNode{
		root: n,
	}
}
