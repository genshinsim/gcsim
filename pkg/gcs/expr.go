package gcs

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func (e *Eval) evalExpr(ex ast.Expr) Obj {
	switch v := ex.(type) {
	case *ast.NumberLit:
		return e.evalNumberLit(v)
	case *ast.StringLit:
		return e.evalStringLit(v)
	case *ast.Ident:
		return e.evalIdent(v)
	case *ast.BinaryExpr:
		return e.evalBinaryExpr(v)
	default:
		return &null{}
	}
}

func (e *Eval) evalNumberLit(n *ast.NumberLit) Obj {
	return &number{
		isInt: n.IsInt,
		ival:  n.IntVal,
		fval:  n.FloatVal,
	}
}

func (e *Eval) evalStringLit(n *ast.StringLit) Obj {
	return &strval{
		str: n.Value,
	}
}

func (e *Eval) evalIdent(n *ast.Ident) Obj {
	//TODO: this should be a variable
	return &null{}
}

func (e *Eval) evalCallExpr(c *ast.CallExpr) Obj {
	//check if it's a system function
	//otherwise check the function map

	//c.Fun is an expression. It needs to evaluate to a FnStmt

	return &null{}
}

func (e *Eval) evalBinaryExpr(b *ast.BinaryExpr) Obj {
	return &null{}
}
