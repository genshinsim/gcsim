package gcs

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) evalStmt(s ast.Stmt) Obj {
	switch v := s.(type) {
	case *ast.BlockStmt:
		return e.evalBlock(v)
	case *ast.LetStmt:
		return e.evalLet(v)
	case *ast.ActionStmt:
		return e.evalAction(v)
	case *ast.ReturnStmt:
		return e.evalReturnStmt(v)
	case *ast.CtrlStmt:
		return e.evalCtrlStmt(v)
	default:
		return &null{}
	}
}

func (e *Eval) evalBlock(b *ast.BlockStmt) Obj {
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

func (e *Eval) evalLet(l *ast.LetStmt) Obj {
	//variable assignment, expr should evaluate to a number
	res := e.evalExpr(l.Val)
	//res should be a number
	v, ok := res.(*number)
	if !ok {
		panic("return does not eval to a number")
	}
	_, exist := e.varMap[l.Ident]
	if exist {
		panic(fmt.Sprintf("variable %v already exists; cannot redeclare", l.Ident.Val))
	}
	e.varMap[l.Ident] = v
	return &null{}
}

func (e *Eval) evalFnStmt(l *ast.FnStmt) Obj {
	_, exist := e.fnMap[l.FunVal]
	if exist {
		panic(fmt.Sprintf("function %v already exists; cannot redeclare", l.FunVal.Val))
	}
	e.fnMap[l.FunVal] = l
	return &null{}
}

func (e *Eval) evalAction(a *ast.ActionStmt) Obj {
	//TODO: should we make a copy of action here??
	e.Work <- *a
	//block until sim is done with the action
	<-e.Next
	return &null{}
}

func (e *Eval) evalReturnStmt(r *ast.ReturnStmt) Obj {
	res := e.evalExpr(r.Val)
	//res should be a number
	if _, ok := res.(*number); !ok {
		panic("return does not eval to a number")
	}
	return &retval{
		res: res,
	}
}

func (e *Eval) evalCtrlStmt(r *ast.CtrlStmt) Obj {
	return &ctrl{
		typ: r.Typ,
	}
}
