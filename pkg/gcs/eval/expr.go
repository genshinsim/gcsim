package eval

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/conditional"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/constant"
)

func (e *Eval) evalExpr(ex ast.Expr, env *Env) (Obj, error) {
	switch v := ex.(type) {
	case *ast.NumberLit:
		return e.evalNumberLit(v), nil
	case *ast.StringLit:
		return e.evalStringLit(v), nil
	case *ast.FuncExpr:
		return e.evalFuncExpr(v, env), nil
	case *ast.Ident:
		return e.evalIdent(v, env)
	case *ast.UnaryExpr:
		return e.evalUnaryExpr(v, env)
	case *ast.BinaryExpr:
		return e.evalBinaryExpr(v, env)
	case *ast.CallExpr:
		return e.evalCallExpr(v, env)
	case *ast.Field:
		return e.evalField(v)
	case *ast.MapExpr:
		return e.evalMap(v, env)
	default:
		return &null{}, nil
	}
}

func (e *Eval) evalNumberLit(n *ast.NumberLit) Obj {
	return &number{
		isFloat: n.IsFloat,
		ival:    n.IntVal,
		fval:    n.FloatVal,
	}
}

func (e *Eval) evalStringLit(n *ast.StringLit) Obj {
	// strip the ""
	return &strval{
		str: strings.Trim(n.Value, "\""),
	}
}

func (e *Eval) evalFuncExpr(n *ast.FuncExpr, env *Env) Obj {
	return &funcval{
		Args:      n.Func.Args,
		Body:      n.Func.Body,
		Signature: n.Func.Signature,
		Env:       NewEnv(env),
	}
}

func (e *Eval) evalIdent(n *ast.Ident, env *Env) (Obj, error) {
	// TODO: this should be a variable
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
	// check number of param matches
	if len(c.Args) != len(fn.Args) {
		return nil, fmt.Errorf("unmatched number of params for fn %v", c.Fun.String())
	}
	// params are just variables assigned to a local env
	// created from the function's env
	local := NewEnv(fn.Env)
	for i, v := range fn.Args {
		param, err := e.evalExpr(c.Args[i], env)
		if err != nil {
			return nil, err
		}
		local.varMap[v.Value] = &param
	}
	res, err := e.evalBlock(fn.Body, local)
	if err != nil {
		return nil, err
	}
	switch v := res.(type) {
	case *retval:
		return v.res, nil
	case *null:
		if _, ok := fn.Signature.ResultType.(*ast.NumberType); ok {
			// force return to 0
			return &number{}, nil
		}
		return &null{}, nil
	case *number:
		if _, ok := fn.Signature.ResultType.(*ast.NumberType); ok {
			// force return to 0
			return v, nil
		}
		return nil, fmt.Errorf("fn %v returned an invalid type; got %v", c.Fun.String(), res.Inspect())
	default:
		return nil, fmt.Errorf("fn %v returned an invalid type; got %v", c.Fun.String(), res.Inspect())
	}
}

func (e *Eval) evalUnaryExpr(b *ast.UnaryExpr, env *Env) (Obj, error) {
	right, err := e.evalExpr(b.Right, env)
	if err != nil {
		return nil, err
	}

	rconst := makeConstant(right)
	if rconst == nil {
		return nil, fmt.Errorf("unary expression does not evaluate to a comparable variable, got %v ", right.Inspect())
	}
	result, err := constant.UnaryOp(b.Op, rconst)
	if err != nil {
		return nil, err
	}
	return fromConstant(result), nil
}

func (e *Eval) evalBinaryExpr(b *ast.BinaryExpr, env *Env) (Obj, error) {
	// eval left, right, operator
	left, err := e.evalExpr(b.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := e.evalExpr(b.Right, env)
	if err != nil {
		return nil, err
	}

	lconst := makeConstant(left)
	if lconst == nil {
		return nil, fmt.Errorf("binary expression does not evaluate to a comparable variable, got %v ", left.Inspect())
	}
	rconst := makeConstant(right)
	if rconst == nil {
		return nil, fmt.Errorf("binary expression does not evaluate to a comparable variable, got %v ", right.Inspect())
	}
	result, err := constant.BinaryOp(b.Op, lconst, rconst)
	if err != nil {
		return nil, err
	}
	return fromConstant(result), nil
}

func (e *Eval) evalField(n *ast.Field) (Obj, error) {
	r, err := conditional.Eval(e.Core, n.Value)
	if err != nil {
		return nil, err
	}

	num := &number{}
	switch v := r.(type) {
	case bool:
		num = bton(v)
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

func (e *Eval) evalMap(m *ast.MapExpr, env *Env) (Obj, error) {
	if len(m.Fields) == 0 { // empty map
		return &mapval{}, nil
	}

	r := &mapval{
		fields: make(map[string]Obj),
	}
	for k, v := range m.Fields {
		obj, err := e.evalExpr(v, env)
		if err != nil {
			return nil, err
		}
		r.fields[k] = obj
	}
	return r, nil
}

func makeConstant(ex Obj) constant.Value {
	switch v := ex.(type) {
	case *number:
		if v.isFloat {
			return constant.Make(v.fval)
		}
		return constant.Make(v.ival)
	case *strval:
		return constant.Make(v.str)
	default:
		return nil
	}
}

func fromConstant(x constant.Value) Obj {
	switch v := constant.Val(x).(type) {
	case int64:
		return &number{ival: v}
	case float64:
		return &number{fval: v, isFloat: true}
	case string:
		return &strval{str: v}
	default:
		return &null{}
	}
}
