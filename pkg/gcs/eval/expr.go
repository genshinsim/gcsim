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
	case *ast.MapExpr:
		return mapExprEval(v)
	case *ast.BinaryExpr:
		return binaryExprEval(v)
	case *ast.UnaryExpr:
		return unaryExprEval(v)
	case *ast.FuncExpr:
		//FuncExpr is only used for anon funcs, followed after a let stmt
		return funcExprEval(v)
	case *ast.Ident:
		return identLitEval(v)
	case *ast.CallExpr:
		return callExprEval(v)
	case *ast.Field:
		//TODO: fields?
		return nil
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

// funcLitEval is never called directly, but is used by either FuncExpr or FuncStmt, both of which
// are essentially wrappers around FuncLit
func funcLitEval(n *ast.FuncLit) evalNode {
	return &funcval{
		Args:      n.Args,
		Body:      n.Body,
		Signature: n.Signature,
	}
}

func funcExprEval(n *ast.FuncExpr) evalNode {
	return funcLitEval(n.Func)
}
