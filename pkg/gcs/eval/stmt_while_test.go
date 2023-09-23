package eval

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalWhileStmt(t *testing.T) {
	prog := `
	let x = 1;
	while x < 5 {
		x = x + 1;
	}
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
	res, _, err := eval.base.nextAction()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val, ok := res.(*number)
	if !ok {
		t.Errorf("res is not number, got %v", res.Typ())
	}
	if val.ival != 5 {
		t.Errorf("expecting res to be 5, got %v", val.ival)
	}
}

func TestEvalWhileWithIfBreakStmt(t *testing.T) {
	prog := `
	let x = 1;
	while 1 {
		x = x + 1;
		if x >= 5 {
			break;
		}
	}
	x;
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
	res, _, err := eval.base.nextAction()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val, ok := res.(*number)
	if !ok {
		t.Errorf("res is not number, got %v", res.Typ())
	}
	if val.ival != 5 {
		t.Errorf("expecting res to be 5, got %v", val.ival)
	}
}

func TestEvalWhileWithIfContinueStmt(t *testing.T) {
	prog := `
	let x = 1;
	let y = 0;
	while 1 {
		x = x + 1;
		if x >= 5 {
			break;
		}
		if x < 2 {
			continue;
		}
		y = y + 1;
	}
	y;
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
	res, _, err := eval.base.nextAction()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val, ok := res.(*number)
	if !ok {
		t.Errorf("res is not number, got %v", res.Typ())
	}
	if val.ival != 3 {
		t.Errorf("expecting res to be 3, got %v", val.ival)
	}
}
