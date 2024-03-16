package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalAssignStmt(t *testing.T) {
	n := &ast.AssignStmt{
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
	var start Obj = &number{ival: 4}
	env.put("someval", &start)
	val, err := runEvalReturnResWhenDone(evalFromNode(n, env))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok := val.(*number)
	if !ok {
		t.Errorf("res is not number, got %v", val.Typ())
	}
	obj, ok := env.v("someval")
	if !ok {
		t.Errorf("someval does not exist in env")
	}
	v, ok := (*obj).(*number)
	if !ok {
		t.Errorf("obj from env is not number, got %v", val.Typ())
	}
	if v.ival != 5 {
		t.Errorf("expected result to be %v, got %v", 5, v.ival)
	}
}
