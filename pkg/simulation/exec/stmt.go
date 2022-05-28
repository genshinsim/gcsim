package exec

import "github.com/genshinsim/gcsim/pkg/parse"

func (e *Executor) evalStmt(s parse.Stmt) Obj {
	switch v := s.(type) {
	case *parse.BlockStmt:
		return e.evalBlock(v)
	case *parse.ActionStmt:
		return e.evalAction(v)
	case *parse.ReturnStmt:
		return e.evalReturnStmt(v)
	case *parse.CtrlStmt:
		return e.evalCtrlStmt(v)
	default:
		return &null{}
	}
}

func (e *Executor) evalBlock(b *parse.BlockStmt) Obj {
	//blocks are effectively a list of statements, so we just need to loop through
	//and evalNode
	for _, n := range b.List {
		v := e.evalNode(n)
		switch v.(type) {
		case *retval:
			// these object should stop execution of current block
			return v
		case *ctrl:
			// TODO: how do we check for invalid continue or break here
			return v
		}
	}
	return &null{}
}

func (e *Executor) evalAction(a *parse.ActionStmt) Obj {
	//TODO: should we make a copy of action here??
	e.Work <- *a
	//block until sim is done with the action
	<-e.Next
	return &null{}
}

func (e *Executor) evalReturnStmt(r *parse.ReturnStmt) Obj {
	res := e.evalExpr(r.Val)
	return &retval{
		res: res,
	}
}

func (e *Executor) evalCtrlStmt(r *parse.CtrlStmt) Obj {
	return &ctrl{
		typ: r.Typ,
	}
}
