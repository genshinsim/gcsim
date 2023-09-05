package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalBasicCallExpr(t *testing.T) {
	n := &ast.CallExpr{
		Fun: &ast.Ident{
			Value: "print",
		},
		Args: []ast.Expr{
			&ast.StringLit{
				Value: "hi",
			},
		},
	}
	eval, err := NewEvaluator(n)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	eval.initSysFuncs(eval.env)
	val, err := runEvaluatorReturnResWhenDone(eval)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok := val.(*null)
	if !ok {
		t.Errorf("res is not null, got %v", val.Typ())
	}
}
