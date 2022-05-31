package gcs

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) print(c *ast.CallExpr, env *Env) Obj {
	//concat all args
	for _, arg := range c.Args {
		val := e.evalExpr(arg, env)
		fmt.Print(val.Inspect())
	}
	fmt.Print("\n")
	return &number{}
}

func (e *Eval) wait(c *ast.CallExpr, env *Env) Obj {
	e.Work <- &ast.ActionStmt{
		Action: action.ActionWait,
		Param:  map[string]int{},
	}

	return &null{}
}
