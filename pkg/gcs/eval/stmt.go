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
		return ctrlStmtEval(v)
	case *ast.IfStmt:
		return ifStmtEval(v, env)
	case *ast.WhileStmt:
		return whileStmtEval(v, env)
	case *ast.ForStmt:
		//TODO: for stmt
		return nil
	case *ast.AssignStmt:
		return assignStmtEval(v, env)
	case *ast.SwitchStmt:
		//TODO: switch
		return nil
	default:
		return nil
	}
}

func ctrlStmtEval(n *ast.CtrlStmt) evalNode {
	return &ctrl{
		typ: n.Typ,
	}
}
