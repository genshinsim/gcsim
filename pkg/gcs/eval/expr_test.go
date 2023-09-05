package eval

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalNum(t *testing.T) {
	n := &ast.NumberLit{
		IntVal: 5,
	}
	e := evalFromNode(n)
	if e == nil {
		t.Error("invalid executor from number")
		t.FailNow()
	}
	val, done, err := e.nextAction(nil)
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("expected node to be done, got false")
	}
	v, ok := val.(*number)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	if v.ival != n.IntVal {
		t.Errorf("expected result to be %v, got %v", n.IntVal, v.ival)
	}
}

func BenchmarkEvalNum(b *testing.B) {
	x := &ast.NumberLit{
		IntVal: 5,
	}
	for n := 0; n < b.N; n++ {
		e := evalFromNode(x)
		e.nextAction(nil)
	}
}

func TestEvalString(t *testing.T) {
	n := &ast.StringLit{
		Value: "bob",
	}
	e := evalFromNode(n)
	if e == nil {
		t.Error("invalid executor from string")
		t.FailNow()
	}
	val, done, err := e.nextAction(nil)
	if err != nil {
		t.Error(err)
	}
	if !done {
		t.Error("expected node to be done, got false")
	}
	v, ok := val.(*strval)
	if !ok {
		t.Errorf("res is not a number, got %v", val.Typ())
	}
	if v.str != n.Value {
		t.Errorf("expected result to be %v, got %v", n.Value, v.str)
	}
}

func TestEvalFuncExpr(t *testing.T) {
	n := &ast.FuncExpr{
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
	val, err := runEvalReturnResWhenDone(evalFromNode(n), nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok := val.(*funcval)
	if !ok {
		t.Errorf("res is not a function, got %v", val.Typ())
	}
}
