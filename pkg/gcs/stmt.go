package gcs

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) evalStmt(s ast.Stmt, env *Env) (Obj, error) {
	switch v := s.(type) {
	case *ast.BlockStmt:
		return e.evalBlock(v, env)
	case *ast.IncDecStmt:
		return e.evalIncDec(v, env)
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
	case *ast.ForStmt:
		return e.evalForStmt(v, env)
	case *ast.AssignStmt:
		return e.evalAssignStmt(v, env)
	case *ast.SwitchStmt:
		return e.evalSwitchStmt(v, env)
	default:
		return &null{}, nil
	}
}

func (e *Eval) evalBlock(b *ast.BlockStmt, env *Env) (Obj, error) {
	// blocks are effectively a list of statements, so we just need to loop through
	// and evalNode
	// blocks should create a new environment
	scope := NewEnv(env)
	for _, n := range b.List {
		v, err := e.evalNode(n, scope)
		if err != nil {
			return nil, err
		}
		switch v.(type) {
		case *retval:
			// these object should stop execution of current block
			return v, nil
		case *ctrl:
			// TODO: how do we check for invalid continue or break here
			// prob need to add some sort of context to env
			return v, nil
		}
	}
	return &null{}, nil
}

func (e *Eval) evalIncDec(i *ast.IncDecStmt, env *Env) (Obj, error) {
	val, err := env.v(i.Ident.Val)
	if err != nil {
		return nil, err
	}
	num := &number{
		ival: 1,
		fval: 1,
	}
	switch i.Val.Typ {
	case ast.ItemInc:
		num = add(val, num)
	case ast.ItemDec:
		num = sub(val, num)
	}
	val.fval = num.fval
	val.ival = num.ival
	val.isFloat = num.isFloat
	return val, nil
}

func (e *Eval) evalLet(l *ast.LetStmt, env *Env) (Obj, error) {
	// variable assignment, expr should evaluate to a number
	res, err := e.evalExpr(l.Val, env)
	if err != nil {
		return nil, err
	}
	// res should be a number
	v, ok := res.(*number)
	// e.Log.Printf("let expr: %v, type: %T\n", res, res)
	if !ok {
		return nil, fmt.Errorf("let expression for %v does evaluate to a number, got %v", l.Ident, res.Inspect())
	}
	_, exist := env.varMap[l.Ident.Val]
	if exist {
		return nil, fmt.Errorf("variable %v already exists; cannot redeclare", l.Ident.Val)
	}
	num := *v // value copying
	env.varMap[l.Ident.Val] = &num
	return &null{}, nil
}

func (e *Eval) evalFnStmt(l *ast.FnStmt, env *Env) (Obj, error) {
	_, exist := env.fnMap[l.FunVal.Val]
	if exist {
		return nil, fmt.Errorf("function %v already exists; cannot redeclare", l.FunVal.Val)
	}
	env.fnMap[l.FunVal.Val] = l
	return &null{}, nil
}

func (e *Eval) evalAssignStmt(a *ast.AssignStmt, env *Env) (Obj, error) {
	res, err := e.evalExpr(a.Val, env)
	if err != nil {
		return nil, err
	}
	v, ok := res.(*number)
	// e.Log.Printf("let expr: %v, type: %T\n", res, res)
	if !ok {
		return nil, fmt.Errorf("value assigned to variable %v does evaluate to a number, got %v", a.Ident, res.Inspect())
	}
	n, err := env.v(a.Ident.Val)
	if err != nil {
		return nil, err
	}
	n.fval = v.fval
	n.ival = v.ival
	n.isFloat = v.isFloat

	return n, nil
}

func (e *Eval) execSwap(char keys.Char) (Obj, error) {
	e.Work <- &ast.ActionStmt{
		Char:   char,
		Action: action.ActionSwap,
	}
	_, ok := <-e.Next
	if !ok {
		return nil, ErrTerminated // no more work, shutting down
	}

	return &null{}, nil
}

func (e *Eval) evalAction(a *ast.ActionStmt, env *Env) (Obj, error) {
	// check if character is active, if not then issue a swap action first
	if !e.Core.Player.CharIsActive(a.Char) {
		res, err := e.execSwap(a.Char)
		if err != nil {
			return nil, err
		}
		if res.Typ() != typNull {
			return res, nil
		}
	}

	// TODO: should we make a copy of action here??
	e.Work <- a
	// block until sim is done with the action; unless we're done
	_, ok := <-e.Next
	if !ok {
		return nil, ErrTerminated // no more work, shutting down
	}
	return &null{}, nil
}

func (e *Eval) evalReturnStmt(r *ast.ReturnStmt, env *Env) (Obj, error) {
	res, err := e.evalExpr(r.Val, env)
	if err != nil {
		return nil, err
	}
	// e.Log.Printf("return res: %v, type: %T\n", res, res)
	// res should be a number
	if _, ok := res.(*number); !ok {
		return nil, fmt.Errorf("return expression does not evaluate to a number, got %v", res.Inspect())
	}
	return &retval{
		res: res,
	}, nil
}

func (e *Eval) evalCtrlStmt(r *ast.CtrlStmt, env *Env) (Obj, error) {
	return &ctrl{
		typ: r.Typ,
	}, nil
}

func (e *Eval) evalIfStmt(i *ast.IfStmt, env *Env) (Obj, error) {
	cond, err := e.evalExpr(i.Condition, env)
	if err != nil {
		return nil, err
	}
	if otob(cond) {
		return e.evalBlock(i.IfBlock, env)
	} else if i.ElseBlock != nil {
		return e.evalStmt(i.ElseBlock, env)
	}
	return &null{}, nil
}

func (e *Eval) evalWhileStmt(w *ast.WhileStmt, env *Env) (Obj, error) {
	for {
		// if condition is false, break
		cond, err := e.evalExpr(w.Condition, env)
		if err != nil {
			return nil, err
		}
		if !otob(cond) {
			break
		}

		// execute block
		res, err := e.evalBlock(w.WhileBlock, env)
		if err != nil {
			return nil, err
		}

		// if result is a break stmt, stop loo
		if t, ok := res.(*ctrl); ok && t.typ == ast.CtrlBreak {
			break
		}
	}
	return &null{}, nil
}

func (e *Eval) evalForStmt(f *ast.ForStmt, env *Env) (Obj, error) {
	scope := NewEnv(env)
	if f.Init != nil {
		e.evalStmt(f.Init, scope)
	}

	for {
		if f.Cond != nil {
			// if condition is false, break
			cond, err := e.evalExpr(f.Cond, scope)
			if err != nil {
				return nil, err
			}
			if !otob(cond) {
				break
			}
		}

		// execute block
		res, err := e.evalBlock(f.Body, scope)
		if err != nil {
			return nil, err
		}

		// if result is a break stmt, stop loo
		if t, ok := res.(*ctrl); ok && t.typ == ast.CtrlBreak {
			break
		}

		if f.Post != nil {
			e.evalStmt(f.Post, scope)
		}
	}
	return &null{}, nil
}

func (e *Eval) evalSwitchStmt(swt *ast.SwitchStmt, env *Env) (Obj, error) {
	cond, err := e.evalExpr(swt.Condition, env)
	if err != nil {
		return nil, err
	}

	// condition should be a number
	// res should be a number
	var v *number = nil
	if _, ok := cond.(*null); !ok {
		val, ok := cond.(*number)
		// e.Log.Printf("let expr: %v, type: %T\n", res, res)
		if !ok {
			return nil, fmt.Errorf("switch condition does not evaluate to a number, got %v", cond.Inspect())
		}
		v = val
	}
	ft := false
	found := false
	// loop through the cases, executing first one that evals true
	for i := range swt.Cases {
		// each case expr needs to evaluate to a number
		cc, err := e.evalExpr(swt.Cases[i].Condition, env)
		if err != nil {
			return nil, err
		}
		c, ok := cc.(*number)
		if !ok {
			return nil, fmt.Errorf("switch case condition does not evaluate to a number, got %v", cc.Inspect())
		}
		if (v == nil && ntob(c)) || (v != nil && ntob(eq(c, v))) || ft {
			found = true
			res, err := e.evalBlock(swt.Cases[i].Body, env)
			if err != nil {
				return nil, err
			}
			e.Log.Printf("res from case block: %v typ %T\n", res, res)
			switch t := res.(type) {
			case *ctrl:
				switch t.typ {
					case ast.CtrlFallthrough:
						ft = true
					case ast.CtrlBreak:
						return &null{}, nil
				}
			default:
				// switch is done
				return res, nil
			}
		}
	}
	if !found || ft {
		if swt.Default == nil {
			return &null{}, nil
		}
		return e.evalBlock(swt.Default, env)
	}
	return &null{}, nil
}
