package eval

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/parser"
)

func TestType(t *testing.T) {
	p := parser.New("type(1);")
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
		eval.Continue()
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
		t.Errorf("expecting type to return string, got %v", result.Typ())
	}
	if eval.Err() != nil {
		t.Error(eval.Err())
	}
}

func TestForceTerminate(t *testing.T) {
	// test terminate eval early should gracefully exit
	p := parser.New(`
	for let i = 0; i < 50; i = i + 1 {
		delay(1);
	}`)
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
	for range 4 {
		eval.Continue()
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
	err = eval.Err()
	if err != nil {
		t.Error(err)
	}
	// confirm that NextAction now returns nil
	for range 4 {
		eval.Continue()
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
	// make sure sleep is evaluated as wait
	p := parser.New("sleep(1);")
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
	eval.Continue()
	a, err := eval.NextAction()
	if err != nil {
		t.Errorf("unexpected error getting next action: %v", err)
	}
	if a == nil {
		t.Error("unexpected next action is nil")
		t.FailNow()
	}
	if a.Action != action.ActionWait {
		t.Errorf("expecting action to be wait, got %v", a.Action.String())
	}
	err = eval.Exit()
	if err != nil {
		t.Errorf("unexpected error exiting: %v", err)
	}
	err = eval.Err()
	if err != nil {
		t.Error(err)
	}
}

func TestDoneCheck(t *testing.T) {
	// eval should exit once out of action; NextAction() should return nil
	p := parser.New(`
	for let i = 0; i < 4; i = i + 1 {
		delay(1);
	}`)
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
		eval.Continue()
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
	// confirm that NextAction continues to return nil
	for range 4 {
		a, err := eval.NextAction()
		if err != nil {
			t.Errorf("unexpected error when checking for NextAction() should be nil: %v", err)
		}
		if a != nil {
			t.Errorf("NextAction() should return nil indicating no more action, got %v", a)
		}
	}
}
