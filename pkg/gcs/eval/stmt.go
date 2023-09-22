package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func evalFromStmt(n ast.Stmt) evalNode {
	switch v := n.(type) {
	case *ast.BlockStmt:
		return blockStmtEval(v)
	case *ast.LetStmt:
		//TODO: let stmt
		return nil
	case *ast.ReturnStmt:
		return returnStmtEval(v)
	case *ast.FnStmt:
		return fnStmtEval(v)
	case *ast.CtrlStmt:
		//TODO: ctrl stmts (break/continue)
		return nil
	case *ast.IfStmt:
		//TODO: if stmts
		return nil
	case *ast.WhileStmt:
		//TODO: while stmt
		return nil
	case *ast.ForStmt:
		//TODO: for stmt
		return nil
	case *ast.AssignStmt:
		//TODO: val assignment
		return nil
	case *ast.SwitchStmt:
		//TODO: switch
		return nil
	default:
		return nil
	}
}
