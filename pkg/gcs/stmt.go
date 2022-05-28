package gcs

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) evalStmt(s ast.Stmt, env *Env) Obj {
	switch v := s.(type) {
	case *ast.BlockStmt:
		return e.evalBlock(v, env)
	case *ast.LetStmt:
		return e.evalLet(v, env)
	case *ast.FnStmt:
		return e.evalFnStmt(v, env)
	case *ast.ActionStmt:
		return e.evalAction(v, env)
	case *ast.ReturnStmt:
		return e.evalReturnStmt(v, env)
	case *ast.CtrlStmt:
		return e.evalCtrlStmt(v, env)
	default:
		return &null{}
	}
}

func (e *Eval) evalBlock(b *ast.BlockStmt, env *Env) Obj {
	//blocks are effectively a list of statements, so we just need to loop through
	//and evalNode
	for _, n := range b.List {
		v := e.evalNode(n, env)
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

func (e *Eval) evalLet(l *ast.LetStmt, env *Env) Obj {
	//variable assignment, expr should evaluate to a number
	res := e.evalExpr(l.Val, env)
	//res should be a number
	v, ok := res.(*number)
	if !ok {
		panic("return does not eval to a number")
	}
	_, exist := env.varMap[l.Ident.Val]
	if exist {
		panic(fmt.Sprintf("variable %v already exists; cannot redeclare", l.Ident.Val))
	}
	env.varMap[l.Ident.Val] = v
	return &null{}
}

func (e *Eval) evalFnStmt(l *ast.FnStmt, env *Env) Obj {
	_, exist := env.fnMap[l.FunVal.Val]
	if exist {
		panic(fmt.Sprintf("function %v already exists; cannot redeclare", l.FunVal.Val))
	}
	env.fnMap[l.FunVal.Val] = l
	return &null{}
}

func (e *Eval) evalAction(a *ast.ActionStmt, env *Env) Obj {
	//TODO: should we make a copy of action here??
	e.Work <- *a
	//block until sim is done with the action
	<-e.Next
	return &null{}
}

func (e *Eval) evalReturnStmt(r *ast.ReturnStmt, env *Env) Obj {
	res := e.evalExpr(r.Val, env)
	//res should be a number
	if _, ok := res.(*number); !ok {
		panic("return does not eval to a number")
	}
	return &retval{
		res: res,
	}
}

func (e *Eval) evalCtrlStmt(r *ast.CtrlStmt, env *Env) Obj {
	return &ctrl{
		typ: r.Typ,
	}
}
