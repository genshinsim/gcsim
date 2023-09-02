package eval

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func evalFromExpr(n ast.Expr) evalNode {
	switch v := n.(type) {
	case *ast.NumberLit:
		return numberLitEval(v)
	case *ast.StringLit:
		return stringLitEval(v)
	case *ast.BinaryExpr:
		return binaryExprEval(v)
	case *ast.UnaryExpr:
		return unaryExprEval(v)
	default:
		return nil

	}
}

func numberLitEval(n *ast.NumberLit) evalNode {
	return &number{
		isFloat: n.IsFloat,
		ival:    n.IntVal,
		fval:    n.FloatVal,
	}
}

func stringLitEval(n *ast.StringLit) evalNode {
	return &strval{
		str: strings.Trim(n.Value, "\""),
	}
}
