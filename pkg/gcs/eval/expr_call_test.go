package eval

import (
	"fmt"
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
	eval, err := NewEvaluator(n, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val, _, err := eval.base.nextAction(eval.env)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok := val.(*null)
	if !ok {
		t.Errorf("res is not null, got %v", val.Typ())
	}
}

func TestEvalFnCall(t *testing.T) {
	prog := `
	fn somefn(a number) number {
		xingqiu attack;
		return a + 1;
	}
	print(somefn(1));
	`
	p := ast.New(prog)
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("program:")
	fmt.Println(gcsl.String())
	eval, err := NewEvaluator(gcsl, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val, _, err := eval.base.nextAction(eval.env)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//expecting the first call to end with an action
	_, ok := val.(*actionval)
	if !ok {
		t.Errorf("res is not an action val, got %v", val.Typ())
	}
	//next call should get end result
	val, _, err = eval.base.nextAction(eval.env)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, ok = val.(*null)
	if !ok {
		t.Errorf("res is not null, got %v", val.Typ())
	}
}
