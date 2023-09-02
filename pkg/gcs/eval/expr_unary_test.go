package eval

import (
	"fmt"
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

	e := evalFromNode(n)
	if e == nil {
		t.Error("invalid executor from number")
		t.FailNow()
	}
	var val Obj
	var done bool
	var err error
	for !done {
		val, done, err = e.evalNext(nil)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(val)
	}
	if !done {
		t.Error("expected node to be done, got false")
	}
	v, ok := val.(*number)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	if v.ival != -4 {
		t.Errorf("expected result to be %v, got %v", -4, v.ival)
	}

}
