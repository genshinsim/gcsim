package ast

import (
	"go/token"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

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

// node() implementations for expression/types nodes
func (*BinaryExpr) Copy()
func (*CompExpr) Copy()
func (*ActionExpr) Copy()

// exprNode()
func (*BinaryExpr) exprNode()
func (*CompExpr) exprNode()
func (*ActionExpr) exprNode()
