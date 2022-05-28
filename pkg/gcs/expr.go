package gcs

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (e *Eval) evalExpr(ex ast.Expr, env *Env) Obj {
	switch v := ex.(type) {
	case *ast.NumberLit:
		return e.evalNumberLit(v, env)
	case *ast.StringLit:
		return e.evalStringLit(v, env)
	case *ast.Ident:
		return e.evalIdent(v, env)
	case *ast.BinaryExpr:
		return e.evalBinaryExpr(v, env)
	default:
		return &null{}
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
	return &strval{
		str: n.Value,
	}
}

func (e *Eval) evalIdent(n *ast.Ident, env *Env) Obj {
	//TODO: this should be a variable
	return env.v(n.Value)
}

func (e *Eval) evalCallExpr(c *ast.CallExpr, env *Env) Obj {
	//check if it's a system function
	//otherwise check the function map

	//c.Fun is an expression. It needs to evaluate to a FnStmt

	return &null{}
}

func (e *Eval) evalBinaryExpr(b *ast.BinaryExpr, env *Env) Obj {
	//eval left, right, operator
	left := e.evalExpr(b.Left, env)
	right := e.evalExpr(b.Right, env)
	//binary expressions should only result in number results
	//otherwise panic for now?
	l, ok := left.(*number)
	if !ok {
		panic(fmt.Sprintf("expr does not evaluate to a number: %v\n", b.Left.String()))
	}
	r, ok := right.(*number)
	if !ok {
		panic(fmt.Sprintf("expr does not evaluate to a number: %v\n", b.Right.String()))
	}
	switch b.Op.Typ {
	case ast.LogicAnd:
		return and(l, r)
	case ast.LogicOr:
		return or(l, r)
	case ast.ItemPlus:
		return add(l, r)
	case ast.ItemMinus:
		return sub(l, r)
	case ast.ItemAsterisk:
		return mul(l, r)
	case ast.ItemForwardSlash:
		return div(l, r)
	case ast.OpGreaterThan:
		return gt(l, r)
	case ast.OpGreaterThanOrEqual:
		return gte(l, r)
	case ast.OpEqual:
		return eq(l, r)
	case ast.OpLessThan:
		return lt(l, r)
	case ast.OpLessThanOrEqual:
		return lte(l, r)
	}
	return &null{}
}
