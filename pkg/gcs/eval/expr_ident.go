package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func identLitEval(n *ast.Ident) evalNode {
	return &identExprEvalNode{
		root: n,
	}
}

type identExprEvalNode struct {
	root *ast.Ident
}

func (i *identExprEvalNode) nextAction(env *Env) (Obj, bool, error) {
	res, err := env.v(i.root.Value)
	if err != nil {
		return nil, false, err
	}

	return *res, true, nil
}
