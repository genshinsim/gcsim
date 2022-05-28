package gcs

import (
	"fmt"

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
