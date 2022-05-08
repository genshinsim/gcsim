package exec

import (
	"go/token"

	"github.com/genshinsim/gcsim/pkg/ast"
	"github.com/genshinsim/gcsim/pkg/core"
)

type Executor struct {
	core       *core.Core
	ast        ast.AST
	ActionChan chan ast.ActionExpr
	Next       chan bool
}

//Exec the next Node. Stop if there's nothing left
func (e *Executor) Run() {
	//this should run until it hits a Action
	//it will then pass the action on a resp channel
	//it will then wait for Next before running again
	return
}

//recursively execute the nodes
func execNode(n ast.Node) {

}

func execStmts(stmts []ast.Stmt) {
	for _, v := range stmts {
		execNode(v)
	}
}

func execFor(n ast.ForStmt) {
	//while condition is true, repeat body block
	for evalCondition(n.Cond) {
		execStmts(n.Body)
	}
}

func evalCondition(c *ast.CondExpr) bool {
	if c.IsLeaf {
		return evalComp(c.Comp)
	}
	left := evalCondition(c.Left)
	right := evalCondition(c.Right)

	switch c.Op {
	case token.OR:
		return left || right
	case token.AND:
		return left && right
	default:
		//urecognized operation in comparison
		return false
	}
}

func evalComp(c *ast.CompExpr) bool {
	return false
}
