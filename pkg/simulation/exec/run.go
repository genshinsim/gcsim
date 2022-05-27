package exec

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
)

type Executor struct {
	Core *core.Core
	AST  parse.Node
	Next chan bool
	Work chan parse.ActionStmt
}

//Run will execute the provided AST. Any genshin specific actions will be passed
//back to the
func (e *Executor) Run() {
	//this should run until it hits an Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	e.evalNode(e.AST)
}

func (e *Executor) evalNode(n parse.Node) {
	switch v := n.(type) {
	case parse.Expr:
		e.evalExpr(v)
	case parse.Stmt:
		e.evalStmt(v)
	}
}
