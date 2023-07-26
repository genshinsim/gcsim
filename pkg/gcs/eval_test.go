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
	res, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	simActions := make(chan *action.ActionEval)
	next := make(chan bool)
	go handleSimActions(simActions, next)
	eval := Eval{
		AST:  res.Program,
		Next: next,
		Work: simActions,
		Log:  log.Default(),
	}
	result := eval.Run()
	if result.Typ() != typStr {
		t.Errorf("expecting type to return string, got %v", typStrings[result.Typ()])
	}

}

func handleSimActions(in chan *action.ActionEval, next chan bool) {
	for {
		next <- true
		x, ok := <-in
		if !ok {
			return
		}
		fmt.Printf("\tExecuting: %v\n", x.Action.String())
	}
}
