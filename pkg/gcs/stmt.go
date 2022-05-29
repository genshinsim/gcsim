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
	case *ast.IfStmt:
		return e.evalIfStmt(v, env)
	case *ast.WhileStmt:
		return e.evalWhileStmt(v, env)
	case *ast.AssignStmt:
		return e.evalAssignStmt(v, env)
	case *ast.SwitchStmt:
		return e.evalSwitchStmt(v, env)
	default:
		return &null{}
	}
}

func (e *Eval) evalBlock(b *ast.BlockStmt, env *Env) Obj {
	//blocks are effectively a list of statements, so we just need to loop through
	//and evalNode
	//blocks should create a new environment
	scope := NewEnv(env)
	for _, n := range b.List {
		v := e.evalNode(n, scope)
		switch v.(type) {
		case *retval:
			// these object should stop execution of current block
			return v
		case *ctrl:
			// TODO: how do we check for invalid continue or break here
			// prob need to add some sort of context to env
			return v
		case *terminate:
			return v //program needs to exit now
		}
	}
	return &null{}
}

func (e *Eval) evalLet(l *ast.LetStmt, env *Env) Obj {
	//variable assignment, expr should evaluate to a number
	res := e.evalExpr(l.Val, env)
	//res should be a number
	v, ok := res.(*number)
	// e.Log.Printf("let expr: %v, type: %T\n", res, res)
	if !ok {
		panic("let expr does not eval to a number")
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

func (e *Eval) evalAssignStmt(a *ast.AssignStmt, env *Env) Obj {
	res := e.evalExpr(a.Val, env)
	v, ok := res.(*number)
	// e.Log.Printf("let expr: %v, type: %T\n", res, res)
	if !ok {
		panic("let expr does not eval to a number")
	}
	n := env.v(a.Ident.Val)
	n.fval = v.fval
	n.ival = v.ival
	n.isFloat = v.isFloat

	return n
}

func (e *Eval) evalAction(a *ast.ActionStmt, env *Env) Obj {
	//TODO: should we make a copy of action here??
	e.Work <- *a
	//block until sim is done with the action; unless we're done
	for {
		select {
		case <-e.Next:
			return &null{}
		case <-e.ctx.Done():
			return &terminate{}
		}
	}
}

func (e *Eval) evalReturnStmt(r *ast.ReturnStmt, env *Env) Obj {
	res := e.evalExpr(r.Val, env)
	// e.Log.Printf("return res: %v, type: %T\n", res, res)
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

func (e *Eval) evalIfStmt(i *ast.IfStmt, env *Env) Obj {
	cond := e.evalExpr(i.Condition, env)
	if otob(cond) {
		return e.evalBlock(i.IfBlock, env)

	} else if i.ElseBlock != nil {
		return e.evalBlock(i.ElseBlock, env)
	}
	return &null{}
}

func (e *Eval) evalWhileStmt(w *ast.WhileStmt, env *Env) Obj {
	for {
		//if condition is false, break
		cond := e.evalExpr(w.Condition, env)
		if !otob(cond) {
			break
		}

		//execute block
		res := e.evalBlock(w.WhileBlock, env)

		//if result is a break stmt, stop loo
		if t, ok := res.(*ctrl); ok && t.typ == ast.CtrlBreak {
			break
		}

		//if terminate then end
		if res.Typ() == typTerminate {
			return res
		}
	}
	return &null{}
}

func (e *Eval) evalSwitchStmt(swt *ast.SwitchStmt, env *Env) Obj {
	cond := e.evalExpr(swt.Condition, env)
	//condition should be a number
	//res should be a number
	v, ok := cond.(*number)
	// e.Log.Printf("let expr: %v, type: %T\n", res, res)
	if !ok {
		panic("switch cond does not eval to a number")
	}
	ft := false
	found := false
	//loop through the cases, executing first one that evals true
	for i := range swt.Cases {
		//each case expr needs to evaluate to a number
		cc := e.evalExpr(swt.Cases[i].Condition, env)
		c, ok := cc.(*number)
		if !ok {
			panic("case expr not a number")
		}
		if ntob(eq(c, v)) || ft {
			found = true
			res := e.evalBlock(swt.Cases[i].Body, env)
			e.Log.Printf("res from case block: %v typ %T\n", res, res)
			switch t := res.(type) {
			case *terminate:
				// terminate if we're done execution
				return t
			case *ctrl:
				// check if fallthrough
				if t.typ == ast.CtrlFallthrough {
					ft = true
				}
			default:
				//switch is done
				return res
			}
		}
	}
	if !found || ft {
		return e.evalBlock(swt.Default, env)
	}
	return &null{}
}
