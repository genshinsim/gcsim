package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalIfStmt(t *testing.T) {
	n := &ast.IfStmt{
		Condition: &ast.BinaryExpr{
			Left: &ast.NumberLit{
				IntVal: 5,
			},
			Right: &ast.NumberLit{
				IntVal: 4,
			},
			Op: ast.Token{
				Typ: ast.OpGreaterThan,
			},
		},
		IfBlock: &ast.BlockStmt{
			List: []ast.Node{
				&ast.NumberLit{IntVal: 10},
				&ast.NumberLit{IntVal: 8},
			},
		},
		ElseBlock: &ast.IfStmt{
			Condition: &ast.StringLit{
				Value: "some string",
			},
			IfBlock: &ast.BlockStmt{
				List: []ast.Node{
					&ast.StringLit{Value: "yay"},
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
	v, ok := val.(*number)
	if !ok {
		t.Errorf("res is not number, got %v", val.Typ().String())
		t.FailNow()
	}
	if v.ival != 8 {
		t.Errorf("expected result to be %v, got %v", 8, v.ival)
	}
	// test else block
	n.Condition = &ast.NumberLit{}
	val, err = runEvalReturnResWhenDone(evalFromNode(n, env))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	vs, ok := val.(*strval)
	if !ok {
		t.Errorf("res is not str, got %v: %v", val.Typ().String(), val.Inspect())
		t.FailNow()
	}
	if vs.str != "yay" {
		t.Errorf("expected result to be %v, got %v", "yay", vs.str)
	}
}
