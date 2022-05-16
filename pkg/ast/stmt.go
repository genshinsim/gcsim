package ast

import "go/token"

type Stmt interface {
	Node
	stmtNode()
}

type (
	// A LetStmt node represents a variable assignment i.e. a = 0
	LetStmt struct {
		LHS string
		RHS []Expr
	}

	ExprStmt struct{}

	//represents a braced statement list
	BlockStmt struct {
		List []Stmt
	}

	// An IfStmt node represents an if statement
	IfStmt struct {
		Cond Expr       // condition
		Body *BlockStmt //block to execute if true
		Else Stmt       // else branch, also
	}

	// A branchStmt node represents a break, continue, goto or fallthrough statement
	BranchStmt struct {
		Tok token.Token
	}

	SwitchStmt struct {
		Body *BlockStmt
	}

	CaseClause struct {
		List []Expr // list of expressions or types; nil means default case
		Body []Stmt
	}

	ForStmt struct {
		Cond *BinaryExpr
		Body *BlockStmt
	}
)

// stmtNode()
