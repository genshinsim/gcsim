package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalBasicMapExpr(t *testing.T) {
	n := &ast.MapExpr{
		Fields: map[string]ast.Expr{
			"test": &ast.NumberLit{IntVal: 5},
		},
	}
	val, err := runEvalReturnResWhenDone(evalFromNode(n), nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v, ok := val.(*mapval)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	res, ok := v.fields["test"]
	if !ok {
		t.Errorf("key test does not exist")
		t.FailNow()
	}
	resv, ok := res.(*number)
	if !ok {
		t.Errorf("res from map is not a number, got %v", res.Typ())
		t.FailNow()
	}
	if resv.ival != 5 {
		t.Errorf("expecting result to be %v, got %v", 5, resv.ival)
	}
}
