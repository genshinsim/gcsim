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
	node()
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

	//A BinaryExpr node represents a binary expression i.e. a > b
	BinaryExpr struct {
		Left  *BinaryExpr
		Right *BinaryExpr
		Op    token.Token
	}

	BasicLit struct {
		Kind  token.Token // token.INT, token.FLOAT, token.CHAR
		Value string      // literal string; eg. 42, 3.14
	}

	//an Ident node represents an identifier
	Ident struct {
		Name string
		Kind IdentKind
	}

	CompExpr struct {
		Field string
		Op    token.Token
		Val   interface{}
	}
)

type IdentKind int

const (
	Bad IdentKind = iota
	Var
	Fun
	Lbl
)

// node() implementations for expression/types nodes
func (*BinaryExpr) node()
func (*CompExpr) node()
func (*ActionExpr) node()

// exprNode()
func (*BinaryExpr) exprNode()
func (*CompExpr) exprNode()
func (*ActionExpr) exprNode()

type Stmt interface {
	Node
	stmtNode()
}

type (
	// An AssingStmt node represents a variable assignment i.e. a = 1
	AssignStmt struct {
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
