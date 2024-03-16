package eval

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestEvalSwitchStmt(t *testing.T) {
	prog := `
	switch 2 {
	case 1:
		100;
	case 2:
		200;
	default:
		300;
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
		t.FailNow()
	}
	if val.ival != 200 {
		t.Errorf("expecting res to be 200, got %v", val.ival)
	}
}

func TestEvalSwitchFallthroughStmt(t *testing.T) {
	prog := `
	switch 2 {
	case 1:
		100;
	case 2:
		200;
		fallthrough;
	default:
		300;
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
		t.FailNow()
	}
	if val.ival != 300 {
		t.Errorf("expecting res to be 300, got %v", val.ival)
	}
}

func TestEvalSwitchNilCondFallthroughStmt(t *testing.T) {
	prog := `
	switch {
	case "":
		1;
	case 1:
		100;
	case 0:
		200;
	default:
		300;
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
		t.FailNow()
	}
	if val.ival != 100 {
		t.Errorf("expecting res to be 300, got %v", val.ival)
	}
}

func TestEvalSwitchNoDefaultStmt(t *testing.T) {
	prog := `
	switch 2 {
	case 1:
		100;
		fallthrough;
	case 0:
		200;
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
	_, ok := res.(*null)
	if !ok {
		t.Errorf("res is not null, got %v", res.Typ())
		t.FailNow()
	}
}
