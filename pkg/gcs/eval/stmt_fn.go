package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

type fnStmtEvalNode struct {
	*ast.FnStmt
	env *Env
}

func (f *fnStmtEvalNode) nextAction() (Obj, bool, error) {
	// functionally, FnStmt is just a special type of let statement
	// add ident to env, then create a new fnval
	_, exist := f.env.varMap[f.Ident.Val]
	if exist {
		return nil, false, fmt.Errorf("function %v already exists; cannot redeclare", f.Ident.Val)
	}
	var fn Obj = &funcval{
		Args:      f.Func.Args,
		Body:      f.Func.Body,
		Signature: f.Func.Signature,
		Env:       NewEnv(f.env),
	}
	f.env.varMap[f.Ident.Val] = &fn
	return &null{}, true, nil
}

func fnStmtEval(n *ast.FnStmt, env *Env) evalNode {
	return &fnStmtEvalNode{FnStmt: n, env: env}
}
