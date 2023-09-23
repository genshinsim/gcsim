package eval

import "github.com/genshinsim/gcsim/pkg/gcs/ast"

func evalFromStmt(n ast.Stmt, env *Env) evalNode {
	switch v := n.(type) {
	case *ast.BlockStmt:
		return blockStmtEval(v, env)
	case *ast.LetStmt:
		return letStmtEval(v, env)
	case *ast.ReturnStmt:
		return returnStmtEval(v, env)
	case *ast.FnStmt:
		return fnStmtEval(v, env)
	case *ast.CtrlStmt:
		//TODO: ctrl stmts (break/continue)
		return nil
	case *ast.IfStmt:
		return ifStmtEval(v, env)
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
