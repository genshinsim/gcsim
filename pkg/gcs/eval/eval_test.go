package eval

import (
	"fmt"
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/parser"
)

func TestType(t *testing.T) {
	file := ast.NewFile()
	p := parser.New(file, "type(1);")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
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
	file := ast.NewFile()
	p := parser.New(file, `
	for let i = 0; i < 50; i = i + 1 {
		delay(1);
	}`)
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
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
			return
		}
		if a == nil {
			t.Error("NextAction() should be not be nil")
			t.FailNow()
			return
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
	file := ast.NewFile()
	p := parser.New(file, "sleep(1);")
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
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
		return
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
	file := ast.NewFile()
	p := parser.New(file, `
	for let i = 0; i < 4; i = i + 1 {
		delay(1);
	}`)
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
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

func TestRunSync(t *testing.T) {
	// the program runs on the calling goroutine and hands every action to exec
	file := ast.NewFile()
	p := parser.New(file, `
	for let i = 0; i < 4; i = i + 1 {
		delay(1);
	}
	wait(2);`)
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
	eval.Log = log.Default()
	var got []action.Action
	err = eval.RunSync(func(a *action.Eval) bool {
		got = append(got, a.Action)
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []action.Action{action.ActionDelay, action.ActionDelay, action.ActionDelay, action.ActionDelay, action.ActionWait}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("expecting actions %v, got %v", want, got)
	}
	if eval.Err() != nil {
		t.Error(eval.Err())
	}
}

func TestRunSyncTerminate(t *testing.T) {
	// exec returning false stops the program right away, even inside a for loop's post statement
	file := ast.NewFile()
	p := parser.New(file, `
	fn next(x) {
		delay(1);
		return x + 1;
	}
	for let i = 0; i < 50; i = next(i) {
		wait(1);
	}`)
	_, gcsl, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	eval, _ := NewEvaluator(file, gcsl, nil)
	eval.Log = log.Default()
	count := 0
	stopped := false
	err = eval.RunSync(func(a *action.Eval) bool {
		count++
		// stop on the first delay, which is in the post statement
		stopped = stopped || a.Action == action.ActionDelay
		return !stopped
	})
	if err != nil {
		t.Errorf("unexpected error from RunSync: %v", err)
	}
	if count != 2 {
		t.Errorf("expecting exec to be called 2 times, got %v", count)
	}
	if eval.Err() != nil {
		t.Error(eval.Err())
	}
}
