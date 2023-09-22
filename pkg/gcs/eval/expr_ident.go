package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func identLitEval(n *ast.Ident, env *Env) evalNode {
	return &identExprEvalNode{
		root: n,
		env:  env,
	}
}

type identExprEvalNode struct {
	root *ast.Ident
	env  *Env
}

func (i *identExprEvalNode) nextAction() (Obj, bool, error) {
	res, err := i.env.v(i.root.Value)
	if err != nil {
		return nil, false, err
	}

	return *res, true, nil
}
