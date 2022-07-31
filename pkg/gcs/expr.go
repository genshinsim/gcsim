package gcs

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/conditional"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) evalExpr(ex ast.Expr, env *Env) (Obj, error) {
	switch v := ex.(type) {
	case *ast.NumberLit:
		return e.evalNumberLit(v, env), nil
	case *ast.StringLit:
		return e.evalStringLit(v, env), nil
	case *ast.Ident:
		return e.evalIdent(v, env)
	case *ast.BinaryExpr:
		return e.evalBinaryExpr(v, env)
	case *ast.CallExpr:
		return e.evalCallExpr(v, env)
	case *ast.Field:
		return e.evalField(v, env)
	default:
		return &null{}, nil
	}
}

func (e *Eval) evalNumberLit(n *ast.NumberLit, env *Env) Obj {
	return &number{
		isFloat: n.IsFloat,
		ival:    n.IntVal,
		fval:    n.FloatVal,
	}
}

func (e *Eval) evalStringLit(n *ast.StringLit, env *Env) Obj {
	//strip the ""
	return &strval{
		str: strings.Trim(n.Value, "\""),
	}
}

func (e *Eval) evalIdent(n *ast.Ident, env *Env) (Obj, error) {
	//TODO: this should be a variable
	return env.v(n.Value)
}

func (e *Eval) evalCallExpr(c *ast.CallExpr, env *Env) (Obj, error) {
	//c.Fun should be an Ident; otherwise panic here
	ident, ok := c.Fun.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("invalid function call %v", c.Fun.String())
	}

	//check if it's a system function
	//otherwise check the function map
	switch s := ident.Value; s {
	case "f":
		return e.f()
	case "rand":
		return e.rand()
	case "randnorm":
		return e.randnorm()
	case "print":
		//print outputs
		return e.print(c, env)
	case "wait":
		//execute wait command
		return e.wait(c, env)
	case "set_target_pos":
		return e.setTargetPos(c, env)
	case "set_player_pos":
		return e.setPlayerPos(c, env)
	case "set_default_target":
		return e.setDefaultTarget(c, env)
	case "set_particle_delay":
		return e.setParticleDelay(c, env)
	default:
		//grab the function first
		fn, err := env.fn(s)
		if err != nil {
			return nil, err
		}
		//check number of param matches
		if len(c.Args) != len(fn.Args) {
			return nil, fmt.Errorf("unmatched number of params for fn %v", s)
		}
		//params are just variables assigned to a local env
		local := NewEnv(env)
		for i, v := range fn.Args {
			param, err := e.evalExpr(c.Args[i], env)
			if err != nil {
				return nil, err
			}
			n, ok := param.(*number)
			if !ok {
				return nil, fmt.Errorf("fn %v param %v does not evaluate to a number, got %v", s, v.Value, param.Inspect())
			}
			local.varMap[v.Value] = n
		}
		res, err := e.evalNode(fn.Body, local)
		if err != nil {
			return nil, err
		}
		switch v := res.(type) {
		case *retval:
			return v.res, nil
		case *null:
			return &number{}, nil
		default:
			return nil, fmt.Errorf("fn %v returned an invalid type; expecting a number got %v", s, res.Inspect())
		}
	}
}

func (e *Eval) evalBinaryExpr(b *ast.BinaryExpr, env *Env) (Obj, error) {
	//eval left, right, operator
	left, err := e.evalExpr(b.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := e.evalExpr(b.Right, env)
	if err != nil {
		return nil, err
	}
	//binary expressions should only result in number results
	//otherwise panic for now?
	l, ok := left.(*number)
	if !ok {
		return nil, fmt.Errorf("binary expression does not evaluate to a number, got %v ", left.Inspect())
	}
	r, ok := right.(*number)
	if !ok {
		return nil, fmt.Errorf("binary expression does not evaluate to a number, got %v ", right.Inspect())
	}
	switch b.Op.Typ {
	case ast.LogicAnd:
		return and(l, r), nil
	case ast.LogicOr:
		return or(l, r), nil
	case ast.ItemPlus:
		return add(l, r), nil
	case ast.ItemMinus:
		return sub(l, r), nil
	case ast.ItemAsterisk:
		return mul(l, r), nil
	case ast.ItemForwardSlash:
		return div(l, r), nil
	case ast.OpGreaterThan:
		return gt(l, r), nil
	case ast.OpGreaterThanOrEqual:
		return gte(l, r), nil
	case ast.OpEqual:
		return eq(l, r), nil
	case ast.OpNotEqual:
		return neq(l, r), nil
	case ast.OpLessThan:
		return lt(l, r), nil
	case ast.OpLessThanOrEqual:
		return lte(l, r), nil
	}
	return &null{}, nil
}

func (e *Eval) evalField(n *ast.Field, env *Env) (Obj, error) {
	r := conditional.Eval(e.Core, n.Value)
	return &number{
		ival: r,
	}, nil
}
