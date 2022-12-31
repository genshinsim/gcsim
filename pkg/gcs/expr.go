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
	case *ast.FuncLit:
		return e.evalFuncLit(v, env), nil
	case *ast.Ident:
		return e.evalIdent(v, env)
	case *ast.UnaryExpr:
		return e.evalUnaryExpr(v, env)
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

func (e *Eval) evalFuncLit(n *ast.FuncLit, env *Env) Obj {
	return &funcval{
		Args: n.Args,
		Body: n.Body,
	}
}

func (e *Eval) evalIdent(n *ast.Ident, env *Env) (Obj, error) {
	//TODO: this should be a variable
	res, err := env.v(n.Value)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (e *Eval) evalCallExpr(c *ast.CallExpr, env *Env) (Obj, error) {
	v, err := e.evalExpr(c.Fun, env)
	if err != nil {
		return nil, err
	}
	switch v.(type) {
	case *funcval:
	case *bfuncval:
	default:
		return nil, fmt.Errorf("invalid function call %v", c.Fun.String())
	}

	if bfn, ok := v.(*bfuncval); ok { // is built-in
		return bfn.Body(c, env)
	}

	fn := v.(*funcval)
	//check number of param matches
	if len(c.Args) != len(fn.Args) {
		return nil, fmt.Errorf("unmatched number of params for fn %v", c.Fun.String())
	}
	//params are just variables assigned to a local env
	local := NewEnv(env)
	for i, v := range fn.Args {
		param, err := e.evalExpr(c.Args[i], env)
		if err != nil {
			return nil, err
		}
		local.varMap[v.Value] = &param
	}
	res, err := e.evalNode(fn.Body, local)
	if err != nil {
		return nil, err
	}
	switch v := res.(type) {
	case *retval:
		return v.res, nil
	case *null:
		return &null{}, nil
	default:
		return nil, fmt.Errorf("fn %v returned an invalid type; got %v", c.Fun.String(), res.Inspect())
	}
}

func (e *Eval) evalUnaryExpr(b *ast.UnaryExpr, env *Env) (Obj, error) {
	right, err := e.evalExpr(b.Right, env)
	if err != nil {
		return nil, err
	}
	//unary expressions should only result in number results
	//otherwise panic for now?
	r, ok := right.(*number)
	if !ok {
		return nil, fmt.Errorf("unary expression does not evaluate to a number, got %v ", right.Inspect())
	}
	switch b.Op.Typ {
	case ast.LogicNot:
		return eq(&number{}, r), nil
	case ast.ItemMinus:
		return sub(&number{}, r), nil
	}
	return &null{}, nil
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
	r, err := conditional.Eval(e.Core, n.Value)
	if err != nil {
		return nil, err
	}

	num := &number{}
	switch v := r.(type) {
	case bool:
		if v {
			num.ival = 1
		}
	case int:
		num.ival = int64(v)
	case int64:
		num.ival = v
	case float64:
		num.fval = v
		num.isFloat = true
	default:
		return nil, fmt.Errorf("field condition '.%v' does not evaluate to a number, got %v", strings.Join(n.Value, "."), v)
	}
	return num, nil
}
