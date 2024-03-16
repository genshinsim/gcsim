package eval

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalForStmt(t *testing.T) {
	prog := `
	let j = 0;
	for let i = 0; i < 5; i = i + 1 {
		print(i);
		j = j + 1;
	}
	j;
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
		t.FailNow()
	}
	if val.ival != 5 {
		t.Errorf("expecting res to be 5, got %v", val.ival)
	}
}

func TestEvalForStmtBreak(t *testing.T) {
	prog := `
	let i = 0;
	for i = 0; i < 5; i = i + 1 {
		print(i);
		if i > 1 {
			break;
		}
	}
	i;
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
		t.FailNow()
	}
	if val.ival != 2 {
		t.Errorf("expecting res to be 2, got %v", val.ival)
	}
}

func TestEvalForStmtContinue(t *testing.T) {
	prog := `
	let x = 0;
	for let i = 0; i < 5; i = i + 1 {
		if i < 2 {
			continue;
		}
		x = x + 1;
		print(x);
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
		t.FailNow()
	}
	if val.ival != 3 {
		t.Errorf("expecting res to be 3, got %v", val.ival)
	}
}

func TestEvalForNoInitCondPostStmtBreak(t *testing.T) {
	prog := `
	let i = 0;
	for {
		print(i);
		i = i + 1;
		if i > 1 {
			break;
		}
	}
	i;
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
		t.FailNow()
	}
	if val.ival != 2 {
		t.Errorf("expecting res to be 2, got %v", val.ival)
	}
}
