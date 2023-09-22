package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalFuncStmt(t *testing.T) {
	n := &ast.FnStmt{
		Ident: ast.Token{
			Val: "somefn",
		},
		Func: &ast.FuncLit{
			Signature: nil, //ignoring signature for this test since we're not validating type for this
			Args: []*ast.Ident{
				{
					Value: "a",
				},
				{
					Value: "b",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Node{
					&ast.NumberLit{IntVal: 1},
					&ast.NumberLit{IntVal: 1},
				},
			},
		},
	}
	env := NewEnv(nil)
	val, err := runEvalReturnResWhenDone(evalFromNode(n, env))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok := val.(*null)
	if !ok {
		t.Errorf("res is not null, got %v", val.Typ())
	}
	obj, err := env.v("somefn")
	if err != nil {
		t.Error(err)
	}
	_, ok = (*obj).(*funcval)
	if !ok {
		t.Errorf("obj from env is not funcval, got %v", val.Typ())
	}

}
