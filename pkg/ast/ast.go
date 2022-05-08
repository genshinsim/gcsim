package ast

import (
	"go/token"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

type AST struct {
	Variables map[string]interface{} //map of variables
	Tree      Node
}

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

// An expression is represented by a tree consisting of one or
// more of the following concrete expression nodes
type (
	//A ActionExpr node represents one action, followed by an optional param list
	//TODO: consider making the args expressions as well?
	ActionExpr struct {
		Type  action.Action
		Param map[string]int
	}

	CondExpr struct {
		Left   *CondExpr
		Right  *CondExpr
		Op     token.Token
		IsLeaf bool
		Comp   *CompExpr
	}

	CompExpr struct {
		Field string
		Op    token.Token
		Val   interface{}
	}
)

type Stmt interface {
	Node
	stmtNode()
}

type AssignStmt struct{}
type ExprStmt struct{}

//represents a braced statement list
type BlockStmt struct {
	List []Stmt
}

// An IfStmt node represents an if statement
type IfStmt struct {
	Cond Expr       // condition
	Body *BlockStmt //block to execute if true
	Else Stmt       // else branch, also
}

// A branchStmt node represents a break, continue, goto or fallthrough statement
type BranchStmt struct {
	Tok token.Token
}

type SwitchStmt struct {
	Body *BlockStmt
}

type CaseClause struct {
	List []Expr // list of expressions or types; nil means default case
	Body []Stmt
}

type ForStmt struct {
	Cond *CondExpr
	Body *BlockStmt
}
