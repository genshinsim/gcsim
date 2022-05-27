package exec

import "github.com/genshinsim/gcsim/pkg/parse"

func (e *Executor) evalStmt(s parse.Stmt) {
	switch v := s.(type) {
	case *parse.BlockStmt:
		e.evalBlock(v)
	case *parse.ActionStmt:
		e.evalAction(v)
	}
}
func (e *Executor) evalBlock(b *parse.BlockStmt) {
	//blocks are effectively a list of statements, so we just need to loop through
	//and evalNode
	for _, n := range b.List {
		e.evalNode(n)
	}
}

func (e *Executor) evalAction(a *parse.ActionStmt) {
	//TODO: should we make a copy of action here??
	e.Work <- *a
	//block until sim is done with the action
	<-e.Next
}
