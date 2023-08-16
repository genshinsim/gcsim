package gcs

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func TestType(t *testing.T) {
	p := ast.New("type(1);")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(gcsl, nil)
	eval.Log = log.Default()
	resultChan := make(chan Obj)
	go func() {
		res, err := eval.Run()
		fmt.Printf("done with result: %v, err: %v\n", res, err)
		resultChan <- res
	}()
	for {
		a, err := eval.NextAction()
		if a == nil {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	result := <-resultChan
	if result.Typ() != typStr {
		t.Errorf("expecting type to return string, got %v", typStrings[result.Typ()])
	}
	if eval.Err() != nil {
		t.Error(eval.Err())
	}
}

func TestForceTerminate(t *testing.T) {
	//test terminate eval early should gracefully exit
	p := ast.New("xingqiu attack:50;")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(gcsl, nil)
	eval.Log = log.Default()
	go func() {
		res, err := eval.Run()
		fmt.Printf("done with result: %v, err: %v\n", res, err)
	}()
	for i := 0; i < 4; i++ {
		a, err := eval.NextAction()
		if err != nil {
			t.Errorf("unexpected error when checking for NextAction(): %v", err)
			t.FailNow()
		}
		if a == nil {
			t.Error("NextAction() should be not be nil")
			t.FailNow()
		}
		fmt.Printf("%v %v\n", a.Char.String(), a.Action.String())
	}
	err = eval.Exit()
	if err != nil {
		t.Error(err)
	}
	//confirm that NextAction now returns nil
	for i := 0; i < 4; i++ {
		a, err := eval.NextAction()
		if err != nil {
			t.Errorf("unexpected error when checking for NextAction() should be nil: %v", err)
		}
		if a != nil {
			t.Errorf("NextAction() should return nil indicating no more action, got %v", a)
		}
	}
}

func TestSleepAsWaitAlias(t *testing.T) {
	//make sure sleep is evaluated as wait
	p := ast.New("sleep(1);")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(gcsl, nil)
	eval.Log = log.Default()
	go func() {
		res, err := eval.Run()
		fmt.Printf("done with result: %v, err: %v\n", res, err)
	}()
	a, err := eval.NextAction()
	if err != nil {
		t.Errorf("unexpected error getting next action: %v", err)
	}
	if a == nil {
		t.Error("unexpected next action is nil")
	}
	if a.Action != action.ActionWait {
		t.Errorf("expecting action to be wait, got %v", a.Action.String())
	}
	err = eval.Exit()
	if err != nil {
		t.Errorf("unexpected error exiting: %v", err)
	}
}

func TestDoneCheck(t *testing.T) {
	//eval should exit once out of action; NextAction() should return nil
	p := ast.New("xingqiu attack, attack, skill, burst;")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(gcsl, nil)
	eval.Log = log.Default()
	go func() {
		res, err := eval.Run()
		fmt.Printf("done with result: %v, err: %v\n", res, err)
	}()
	count := 0
	for {
		a, err := eval.NextAction()
		if a == nil {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v %v\n", a.Char.String(), a.Action.String())
		count++
	}
	if count != 4 {
		t.Errorf("expecting NextAction to be called 4 times, got %v", count)
	}
	//confirm that NextAction continues to return nil
	for i := 0; i < 4; i++ {
		a, err := eval.NextAction()
		if err != nil {
			t.Errorf("unexpected error when checking for NextAction() should be nil: %v", err)
		}
		if a != nil {
			t.Errorf("NextAction() should return nil indicating no more action, got %v", a)
		}
	}
}
