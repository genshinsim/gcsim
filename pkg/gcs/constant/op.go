package constant

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func ToBool(v Value) bool {
	switch x := v.(type) {
	case *number:
		return ntob(x)
	case *strval:
		return true
	default:
		return false
	}
}

func ntob(v *number) bool {
	// check int 0
	if !v.isFloat && v.ival == 0 {
		return false
	}
	// check float 0
	if v.isFloat && v.fval == 0 {
		return false
	}
	// otherwise true
	return true
}

func bton(b bool) *number {
	if b {
		return &number{ival: 1, fval: 1}
	}
	return &number{}
}

func and(l, r *number) *number {
	return bton(ntob(l) && ntob(r))
}

func or(l, r *number) *number {
	return bton(ntob(l) || ntob(r))
}

func ntof(v *number) float64 {
	if v.isFloat {
		return v.fval
	}
	return float64(v.ival)
}

func ntoi(v *number) int64 {
	if v.isFloat {
		return int64(v.fval)
	}
	return v.ival
}

func gt(l, r *number) *number {
	return bton(ntof(l) > ntof(r))
}

func gte(l, r *number) *number {
	return bton(ntof(l) >= ntof(r))
}

func lt(l, r *number) *number {
	return bton(ntof(l) < ntof(r))
}

func lte(l, r *number) *number {
	return bton(ntof(l) <= ntof(r))
}

func eq(l, r *number) *number {
	return bton(ntof(l) == ntof(r))
}

func neq(l, r *number) *number {
	return bton(ntof(l) != ntof(r))
}

func add(l, r *number) *number {
	n := &number{
		isFloat: l.isFloat || r.isFloat,
	}
	if n.isFloat {
		n.fval = ntof(l) + ntof(r)
	} else {
		n.ival = l.ival + r.ival
	}
	return n
}

func mul(l, r *number) *number {
	n := &number{
		isFloat: l.isFloat || r.isFloat,
	}
	if n.isFloat {
		n.fval = ntof(l) * ntof(r)
	} else {
		n.ival = l.ival * r.ival
	}
	return n
}

func div(l, r *number) *number {
	n := &number{
		isFloat: l.isFloat || r.isFloat,
	}
	if n.isFloat {
		n.fval = ntof(l) / ntof(r)
	} else {
		n.ival = l.ival / r.ival
	}
	return n
}

func sub(l, r *number) *number {
	n := &number{
		isFloat: l.isFloat || r.isFloat,
	}
	if n.isFloat {
		n.fval = ntof(l) - ntof(r)
	} else {
		n.ival = l.ival - r.ival
	}
	return n
}

func UnaryOp(op ast.Token, right Value) (Value, error) {
	switch op.Typ {
	case ast.ItemPlus:
		if right, ok := right.(*number); ok {
			return right, nil
		}
	case ast.LogicNot:
		return bton(!ToBool(right)), nil
	case ast.ItemMinus:
		if right, ok := right.(*number); ok {
			return sub(&number{}, right), nil
		}
	}

	return nil, fmt.Errorf("invalid unary operator %v%v", op, right.Inspect())
}

func BinaryOp(op ast.Token, left, right Value) (Value, error) {
	if left, ok := left.(*number); ok {
		right, ok := right.(*number)
		if !ok {
			return nil, fmt.Errorf("invalid binary operator %v%v%v", left.Inspect(), op, right.Inspect())
		}

		switch op.Typ {
		case ast.ItemPlus:
			return add(left, right), nil
		case ast.ItemMinus:
			return sub(left, right), nil
		case ast.ItemAsterisk:
			return mul(left, right), nil
		case ast.ItemForwardSlash:
			if !right.isFloat && right.ival == 0 {
				return nil, errors.New("division by zero")
			}
			return div(left, right), nil
		case ast.OpGreaterThan:
			return gt(left, right), nil
		case ast.OpGreaterThanOrEqual:
			return gte(left, right), nil
		case ast.OpEqual:
			return eq(left, right), nil
		case ast.OpNotEqual:
			return neq(left, right), nil
		case ast.OpLessThan:
			return lt(left, right), nil
		case ast.OpLessThanOrEqual:
			return lte(left, right), nil
		case ast.LogicAnd:
			return and(left, right), nil
		case ast.LogicOr:
			return or(left, right), nil
		}
	}

	return nil, fmt.Errorf("invalid binary operator %v%v%v", left.Inspect(), op, right.Inspect())
}
