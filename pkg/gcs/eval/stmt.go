package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func evalFromStmt(n ast.Stmt) evalNode {
	switch v := n.(type) {
	case *ast.BlockStmt:
		return blockStmtEval(v)
	case *ast.ReturnStmt:
		return returnStmtEval(v)
	case *ast.FnStmt:
		return fnStmtEval(v)
	default:
		return nil

	}
}

