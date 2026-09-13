package parser

import (
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/constant"
)

func evalConstant(ex ast.Expr) constant.Value {
	switch v := ex.(type) {
	case *ast.NumberLit:
		if v.IsFloat {
			return constant.Make(v.FloatVal)
		}
		return constant.Make(v.IntVal)
	case *ast.StringLit:
		return constant.Make(v.Value)
	default:
		return nil
	}
}

func constToExpr(pos ast.Pos, x constant.Value) ast.Expr {
	switch val := constant.Val(x).(type) {
	case int64:
		return &ast.NumberLit{
			Pos:     pos,
			IntVal:  val,
			IsFloat: false,
		}
	case float64:
		return &ast.NumberLit{
			Pos:      pos,
			FloatVal: val,
			IsFloat:  true,
		}
	case string:
		return &ast.StringLit{
			Pos:   pos,
			Value: val,
		}
	default:
		return nil
	}
}

func foldConstants(ex ast.Expr) (ast.Expr, error) {
	switch ex := ex.(type) {
	case *ast.UnaryExpr:
		right, err := foldConstants(ex.Right)
		if err != nil {
			return nil, err
		}
		r := evalConstant(right)
		if r == nil {
			return ex, nil
		}
		val, err := constant.UnaryOp(ex.Op, r)
		if err != nil {
			return nil, err
		}
		return constToExpr(ex.Pos, val), nil
	case *ast.BinaryExpr:
		left, err := foldConstants(ex.Left)
		if err != nil {
			return nil, err
		}
		right, err := foldConstants(ex.Right)
		if err != nil {
			return nil, err
		}

		l := evalConstant(left)
		r := evalConstant(right)
		if l == nil || r == nil {
			return ex, nil
		}
		val, err := constant.BinaryOp(ex.Op, l, r)
		if err != nil {
			return nil, err
		}
		return constToExpr(ex.Pos, val), nil
	default:
		return ex, nil
	}
}
