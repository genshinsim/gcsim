package eval

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func evalFromExpr(n ast.Expr, env *Env) evalNode {
	switch v := n.(type) {
	case *ast.NumberLit:
		return numberLitEval(v)
	case *ast.StringLit:
		return stringLitEval(v)
	case *ast.MapExpr:
		return mapExprEval(v, env)
	case *ast.BinaryExpr:
		return binaryExprEval(v, env)
	case *ast.UnaryExpr:
		return unaryExprEval(v, env)
	case *ast.FuncExpr:
		// FuncExpr is only used for anon funcs, followed after a let stmt
		return funcExprEval(v, env)
	case *ast.Ident:
		return identLitEval(v, env)
	case *ast.CallExpr:
		return callExprEval(v, env)
	case *ast.Field:
		return fieldsExprEval(v, env)
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
func funcLitEval(n *ast.FuncLit, env *Env) evalNode {
	return &funcval{
		Args:      n.Args,
		Body:      n.Body,
		Signature: n.Signature,
		Env:       NewEnv(env),
	}
}

func funcExprEval(n *ast.FuncExpr, env *Env) evalNode {
	return funcLitEval(n.Func, env)
}
