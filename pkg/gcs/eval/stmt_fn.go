package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type fnStmtEvalNode struct {
	root *ast.FnStmt
}

func (f *fnStmtEvalNode) evalNext(env *Env) (Obj, bool, error) {
	//add ident to env, then create a new fnval
	_, exist := env.varMap[f.root.Ident.Val]
	if exist {
		return nil, false, fmt.Errorf("function %v already exists; cannot redeclare", f.root.Ident.Val)
	}
	var fn Obj = &funcval{
		Args:      f.root.Func.Args,
		Body:      f.root.Func.Body,
		Signature: f.root.Func.Signature,
		Env:       NewEnv(env),
	}
	env.varMap[f.root.Ident.Val] = &fn
	return &null{}, true, nil

}

func fnStmtEval(n *ast.FnStmt) evalNode {
	//functionally, FnStmt is just a special type of let statement
	return &fnStmtEvalNode{root: n}
}
