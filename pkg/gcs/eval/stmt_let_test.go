package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalLetStmt(t *testing.T) {
	n := &ast.LetStmt{
		Ident: ast.Token{
			Val: "someval",
		},
		Val: &ast.BinaryExpr{
			Left: &ast.BinaryExpr{
				Left: &ast.NumberLit{
					IntVal: 5,
				},
				Right: &ast.NumberLit{
					IntVal: 5,
				},
				Op: ast.Token{
					Typ: ast.ItemPlus,
				},
			},
			Right: &ast.NumberLit{
				IntVal: 5,
			},
			Op: ast.Token{
				Typ: ast.ItemMinus,
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
	obj, err := env.v("someval")
	if err != nil {
		t.Error(err)
	}
	v, ok := (*obj).(*number)
	if !ok {
		t.Errorf("obj from env is not funcval, got %v", val.Typ())
	}
	if v.ival != 5 {
		t.Errorf("expected result to be %v, got %v", 5, v.ival)
	}
}
