package eval

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

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
	res, ok := i.env.v(i.root.Value)
	if !ok {
		return nil, false, fmt.Errorf("ident %v doeas not exist", i.root.Value)
	}

	return *res, true, nil
}
