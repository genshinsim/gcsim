package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalBasicUnaryExpr(t *testing.T) {
	n := &ast.UnaryExpr{
		Right: &ast.NumberLit{
			IntVal: 4,
		},
		Op: ast.Token{
			Typ: ast.ItemMinus,
		},
	}

	val, err := runEvalReturnResWhenDone(evalFromNode(n))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v, ok := val.(*number)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	if v.ival != -4 {
		t.Errorf("expected result to be %v, got %v", -4, v.ival)
	}

}
